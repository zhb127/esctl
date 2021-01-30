package es

import (
	"encoding/base64"
	"esctl/pkg/log"
	"esctl/pkg/util/converttype/tostr"
	"net/http"
	"strings"
	"time"

	goES "github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
)

func newRawClient(config HelperConfig, logHelper log.IHelper) (*goES.Client, error) {
	goESConfig := goES.Config{
		Addresses: strings.Split(config.Addresses, ","),
		Username:  config.Username,
		Password:  config.Password,
	}

	if logHelper != nil {
		goESConfig.Logger = newRawClientLogger(logHelper)
	}

	if config.CertVerify == true {
		if config.CertData == "" {
			return nil, errors.New("config.CertData is empty")
		}

		certByteArr, err := base64.StdEncoding.DecodeString(config.CertData)
		if err != nil {
			return nil, errors.Wrap(err, "config.CertData is invalid")
		}
		goESConfig.CACert = certByteArr
	}

	if config.CertVerify == false {
		httpTransportTmp := http.DefaultTransport
		httpTransport, _ := httpTransportTmp.(*http.Transport)
		httpTransport = httpTransport.Clone()
		httpTransport.TLSClientConfig.InsecureSkipVerify = true
		goESConfig.Transport = httpTransport
	}

	goESClientInst, err := goES.NewClient(goESConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to goES.NewClient")
	}

	// 检查是否可访问
	resp, err := goESClientInst.Info()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get ES Info")
	}

	defer resp.Body.Close()
	if resp.IsError() {
		tmpErr := errors.New(resp.String())
		return nil, errors.Wrap(tmpErr, "failed to get ES Info")
	}

	return goESClientInst, nil
}

func newRawClientLogger(logHelper log.IHelper) *goESClientLogger {
	if logHelper == nil {
		return nil
	}

	l := &goESClientLogger{logHelper: logHelper}
	return l
}

type goESClientLogger struct {
	logHelper log.IHelper
}

func (l *goESClientLogger) RequestBodyEnabled() bool { return true }

func (l *goESClientLogger) ResponseBodyEnabled() bool { return true }

func (l *goESClientLogger) LogRoundTrip(req *http.Request, resp *http.Response, err error, start time.Time, dur time.Duration) error {
	logFields := map[string]interface{}{
		"dur":  dur,
		"req":  nil,
		"resp": nil,
	}

	if err != nil {
		logFields["error"] = err.Error()
	}

	if req != nil {
		reqBodyStr := ""

		// 只记录查询请求体
		reqUrl := req.URL.String()
		if strings.ContainsAny(reqUrl, "/_search") {
			if req.Body != nil && req.Body != http.NoBody {
				reqBodyStr = tostr.FromIOReadCloser(req.Body)
			}
		} else {
			reqBodyStr = "ignore logging request body due to it is not search request"
		}

		logFields["req"] = map[string]interface{}{
			"method": req.Method,
			"url":    req.URL.String(),
			"body":   reqBodyStr,
		}

	}

	if resp != nil {
		respBodyStr := ""
		if resp.Body != nil && resp.Body != http.NoBody {
			// 仅在发生错误时，记录响应体
			if resp.StatusCode >= http.StatusBadRequest {
				respBodyStr = tostr.FromIOReadCloser(resp.Body)
			}
		}

		logFields["resp"] = map[string]interface{}{
			"status": resp.StatusCode,
			"body":   respBodyStr,
		}
	}

	switch {
	case err != nil:
		l.logHelper.Error("an error occurred", logFields)
	case resp != nil && resp.StatusCode < http.StatusBadRequest:
		l.logHelper.Info(resp.Status, logFields)
	case resp != nil && resp.StatusCode >= http.StatusBadRequest && resp.StatusCode < http.StatusInternalServerError:
		l.logHelper.Warn(resp.Status, logFields)
	case resp != nil && resp.StatusCode >= http.StatusInternalServerError:
		l.logHelper.Error(resp.Status, logFields)
	default:
		l.logHelper.Error("no response", logFields)
	}

	return nil
}
