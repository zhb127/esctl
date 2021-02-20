package es

import (
	"bytes"
	"encoding/json"
	"esctl/pkg/log"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	goES "github.com/elastic/go-elasticsearch/v7"
)

type HelperConfig struct {
	Addresses  string `env:"ES_ADDRESSES,required"`
	Username   string `env:"ES_USERNAME"`
	Password   string `env:"ES_PASSWORD"`
	CertData   string `env:"ES_CERT_DATA"`
	CertVerify bool   `env:"ES_CERT_VERIFY,default=false"`
}

type IHelper interface {
	SaveDoc(indexName string, docID string, docBody []byte) error
	DeleteDoc(indexName string, docID string) error
	SearchDocs(indexName string, searchBody []byte) (*SearchDocsResp, error)
	CatIndices(indexNameWildcardExps ...string) (*CatIndicesResp, error)
	CreateIndex(indexName string, indexBody []byte) (*CreateIndexResp, error)
	DeleteIndices(indexNames ...string) (*DeleteIndexResp, error)
	Reindex(srcIndexName string, destIndexName string) (*ReindexResp, error)
	AliasIndex(indexName string, alias string) (*AliasOrUnaliasIndexResp, error)
	UnaliasIndex(indexName string, aliases []string) (*AliasOrUnaliasIndexResp, error)
	ListAliases() (*ListAliasesResp, error)
}

type helper struct {
	config    HelperConfig
	logHelper log.IHelper
	rawClient *goES.Client
}

func NewHelper(config HelperConfig, logHelper log.IHelper) (IHelper, error) {
	// 处理日志
	if logHelper != nil {
		logHelper = logHelper.NewChild().SetWithField("pkg", "es_helper")
	}

	// 处理 go ES 客户端
	rawClient, err := newRawClient(config, logHelper)
	if err != nil {
		return nil, errors.Wrap(err, "failed to newGoESClient")
	}

	h := &helper{
		config,
		logHelper,
		rawClient,
	}

	return h, nil
}

func (h *helper) SaveDoc(indexName string, docID string, docBody []byte) error {
	resp, err := h.rawClient.Index(
		indexName,
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

func (h *helper) DeleteDoc(indexName string, docID string) error {
	resp, err := h.rawClient.Delete(indexName, docID)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		return errors.New(resp.String())
	}

	return nil
}

func (h *helper) SearchDocs(indexName string, searchBody []byte) (*SearchDocsResp, error) {
	resp, err := h.rawClient.Search(
		h.rawClient.Search.WithIndex(indexName),
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

	result := &SearchDocsResp{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

func (h *helper) CatIndices(indexNameWildcardExps ...string) (*CatIndicesResp, error) {
	resp, err := h.rawClient.Cat.Indices(
		h.rawClient.Cat.Indices.WithIndex(indexNameWildcardExps...),
		h.rawClient.Cat.Indices.WithFormat("json"),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}

	res := &CatIndicesResp{}
	if err := json.NewDecoder(resp.Body).Decode(&res.Items); err != nil {
		return nil, err
	}

	return res, nil
}

func (h *helper) CreateIndex(indexName string, indexBody []byte) (*CreateIndexResp, error) {
	resp, err := h.rawClient.Indices.Create(
		indexName,
		h.rawClient.Indices.Create.WithBody(bytes.NewReader(indexBody)),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}

	res := &CreateIndexResp{}
	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		return nil, err
	}

	return res, nil
}

func (h *helper) DeleteIndices(indexNames ...string) (*DeleteIndexResp, error) {
	resp, err := h.rawClient.Indices.Delete(indexNames)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}

	res := &DeleteIndexResp{}
	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		return nil, err
	}

	return res, nil
}

func (h *helper) Reindex(srcIndexName string, destIndexName string) (*ReindexResp, error) {
	body := fmt.Sprintf(`{"source":{"index":"%s"},"dest":{"index":"%s"}}`, srcIndexName, destIndexName)
	resp, err := h.rawClient.Reindex(strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}

	res := &ReindexResp{}
	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		return nil, err
	}

	return res, nil
}

func (h *helper) AliasIndex(indexName string, alias string) (*AliasOrUnaliasIndexResp, error) {
	resp, err := h.rawClient.Indices.PutAlias([]string{indexName}, alias)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}

	res := &AliasOrUnaliasIndexResp{}
	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		return nil, err
	}

	return res, nil
}

func (h *helper) UnaliasIndex(indexName string, aliases []string) (*AliasOrUnaliasIndexResp, error) {
	resp, err := h.rawClient.Indices.DeleteAlias([]string{indexName}, aliases)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}

	res := &AliasOrUnaliasIndexResp{}
	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		return nil, err
	}

	return res, nil
}

func (h *helper) ListAliases() (*ListAliasesResp, error) {
	resp, err := h.rawClient.Cat.Aliases(
		h.rawClient.Cat.Aliases.WithFormat("json"),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}

	res := &ListAliasesResp{}
	if err := json.NewDecoder(resp.Body).Decode(&res.Items); err != nil {
		return nil, err
	}

	return res, nil
}
