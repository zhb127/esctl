package es

import (
	"bytes"
	"encoding/json"
	"esctl/pkg/log"

	goES "github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
)

type HelperConfig struct {
	Addresses  string `env:"ES_ADDRESSES,required"`
	Username   string `env:"ES_USERNAME"`
	Password   string `env:"ES_PASSWORD"`
	CertPath   string `env:"ES_CERT_PATH"`
	CertVerify bool   `env:"ES_CERT_VERIFY,default=false"`
}

type IHelper interface {
	SaveDoc(index string, docID string, docBody []byte) error
	DeleteDoc(index string, docID string) error
	SearchDocs(index string, searchBody []byte) (*SearchDocsResponse, error)
}

type helper struct {
	config     *HelperConfig
	goESClient *goES.Client
}

func NewHelper(config *HelperConfig, logHelper log.IHelper) (*helper, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	goESClient, err := newGoESClient(config, logHelper)
	if err != nil {
		return nil, errors.Wrap(err, "failed to newGoESClient")
	}

	h := &helper{
		config:     config,
		goESClient: goESClient,
	}

	return h, nil
}

func (h *helper) SaveDoc(index string, docID string, docBody []byte) error {
	resp, err := h.goESClient.Index(
		index,
		bytes.NewReader(docBody),
		h.goESClient.Index.WithDocumentID(docID),
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
	resp, err := h.goESClient.Delete(index, docID)
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
	resp, err := h.goESClient.Search(
		h.goESClient.Search.WithIndex(index),
		h.goESClient.Search.WithBody(bytes.NewReader(searchBody)),
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
