package es

import (
	"bytes"
	"encoding/json"
	"esctl/pkg/log"

	"github.com/pkg/errors"

	goES "github.com/elastic/go-elasticsearch/v7"
)

type HelperConfig struct {
	Addresses  string `env:"ES_ADDRESSES,required"`
	Username   string `env:"ES_USERNAME"`
	Password   string `env:"ES_PASSWORD"`
	CertPath   string `env:"ES_CERT_PATH"`
	CertVerify bool   `env:"ES_CERT_VERIFY,default=false"`
}

type IHelper interface {
	Client() *goES.Client
	SaveDoc(index string, docID string, docBody []byte) error
	DeleteDoc(index string, docID string) error
	SearchDocs(index string, searchBody []byte) (*SearchDocsResponse, error)
}

type helper struct {
	config    HelperConfig
	logHelper log.IHelper
	rawClient *goES.Client
}

func NewHelper(config HelperConfig, logHelper log.IHelper) (*helper, error) {

	// 处理日志
	if logHelper != nil {
		logHelper = logHelper.NewChild().SetWithField("pkg", "es_helper")
	}

	// 处理 go ES 客户端
	rawClient, err := newRawClient(config, logHelper)
	if err != nil {
		return nil, errors.Wrap(err, "failed to newGoESClient")
	}

	inst := &helper{
		config,
		logHelper,
		rawClient,
	}

	return inst, nil
}

func (h *helper) Client() *goES.Client {
	return h.rawClient
}

func (h *helper) SaveDoc(index string, docID string, docBody []byte) error {
	resp, err := h.rawClient.Index(
		index,
		bytes.NewReader(docBody),
		h.rawClient.Index.WithDocumentID(docID),
	)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.IsError() {
		return errors.New(resp.String())
	}

	return nil
}

func (h *helper) DeleteDoc(index string, docID string) error {
	resp, err := h.rawClient.Delete(index, docID)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.IsError() {
		return errors.New(resp.String())
	}

	return nil
}

func (h *helper) SearchDocs(index string, searchBody []byte) (*SearchDocsResponse, error) {
	resp, err := h.rawClient.Search(
		h.rawClient.Search.WithIndex(index),
		h.rawClient.Search.WithTrackTotalHits(true),
		h.rawClient.Search.WithBody(bytes.NewReader(searchBody)),
	)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}

	result := &SearchDocsResponse{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}
