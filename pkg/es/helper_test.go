package es

import (
	"esctl/pkg/config/dotenv"
	"esctl/pkg/log"
	"esctl/pkg/util/converttype/tostr"
	tdLog "esctl/test/data/pkg/log"
	"reflect"
	"testing"
	"time"

	goES "github.com/elastic/go-elasticsearch/v7"
	"github.com/stretchr/testify/assert"
)

func MockHelperConfig() HelperConfig {
	config := HelperConfig{}
	if err := dotenv.Decode(&config); err != nil {
		panic(err)
	}
	return config
}

func MockRawClient(config HelperConfig, logHelper log.IHelper) *goES.Client {
	rawClient, err := newRawClient(config, logHelper)
	if err != nil {
		panic(err)
	}
	return rawClient
}

func MockIndexNameExisting() string {
	config := HelperConfig{}
	if err := dotenv.Decode(&config); err != nil {
		panic(err)
	}

	inst, err := NewHelper(config, nil)
	if err != nil {
		panic(err)
	}

	indexName := "mock-index-" + tostr.FromInt64(time.Now().UnixNano())
	if _, err := inst.CreateIndex(indexName, []byte(`{"mappings":{"properties":{"id":{"type":"long"}}}}`)); err != nil {
		panic(err)
	}

	return indexName
}

func TestNewHelper(t *testing.T) {
	type args struct {
		config    HelperConfig
		logHelper log.IHelper
	}
	tests := []struct {
		name    string
		args    args
		want    *helper
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHelper(tt.args.config, tt.args.logHelper)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHelper() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHelper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_helper_SaveDoc(t *testing.T) {
	type fields struct {
		config    HelperConfig
		logHelper log.IHelper
		rawClient *goES.Client
	}
	type args struct {
		index   string
		docID   string
		docBody []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &helper{
				config:    tt.fields.config,
				logHelper: tt.fields.logHelper,
				rawClient: tt.fields.rawClient,
			}
			if err := h.SaveDoc(tt.args.index, tt.args.docID, tt.args.docBody); (err != nil) != tt.wantErr {
				t.Errorf("helper.SaveDoc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_helper_DeleteDoc(t *testing.T) {
	type fields struct {
		config    HelperConfig
		logHelper log.IHelper
		rawClient *goES.Client
	}
	type args struct {
		index string
		docID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &helper{
				config:    tt.fields.config,
				logHelper: tt.fields.logHelper,
				rawClient: tt.fields.rawClient,
			}
			if err := h.DeleteDoc(tt.args.index, tt.args.docID); (err != nil) != tt.wantErr {
				t.Errorf("helper.DeleteDoc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_helper_SearchDocs(t *testing.T) {
	type fields struct {
		config    HelperConfig
		logHelper log.IHelper
		rawClient *goES.Client
	}
	type args struct {
		index    string
		condBody []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *SearchDocsResp
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &helper{
				config:    tt.fields.config,
				logHelper: tt.fields.logHelper,
				rawClient: tt.fields.rawClient,
			}
			got, err := h.SearchDocs(tt.args.index, tt.args.condBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("helper.SearchDocs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("helper.SearchDocs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_helper_CatIndices(t *testing.T) {
	config := MockHelperConfig()
	logHelper := tdLog.MockHelper()
	rawClient := MockRawClient(config, logHelper)

	type fields struct {
		config    HelperConfig
		logHelper log.IHelper
		rawClient *goES.Client
	}
	tests := []struct {
		name    string
		fields  fields
		want    *CatIndicesResp
		wantErr bool
	}{
		{
			fields: fields{
				config:    config,
				logHelper: logHelper,
				rawClient: rawClient,
			},
			want:    &CatIndicesResp{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &helper{
				config:    tt.fields.config,
				logHelper: tt.fields.logHelper,
				rawClient: tt.fields.rawClient,
			}
			got, err := h.CatIndices()
			if (err != nil) != tt.wantErr {
				t.Errorf("helper.ListIndices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
		})
	}
}

func Test_helper_CreateIndex(t *testing.T) {
	config := MockHelperConfig()
	logHelper := tdLog.MockHelper()
	rawClient := MockRawClient(config, logHelper)

	type fields struct {
		config    HelperConfig
		logHelper log.IHelper
		rawClient *goES.Client
	}
	type args struct {
		index     string
		indexBody []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CreateIndexResp
		wantErr bool
	}{
		{
			fields: fields{
				config:    config,
				logHelper: logHelper,
				rawClient: rawClient,
			},
			args: args{
				index:     "test-index",
				indexBody: []byte(`{"mappings":{"properties":{"id":{"type":"long"}}}}`),
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &helper{
				config:    tt.fields.config,
				logHelper: tt.fields.logHelper,
				rawClient: tt.fields.rawClient,
			}
			got, err := h.CreateIndex(tt.args.index, tt.args.indexBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("helper.CreateIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("helper.CreateIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_helper_DeleteIndices(t *testing.T) {
	config := MockHelperConfig()
	logHelper := tdLog.MockHelper()
	rawClient := MockRawClient(config, logHelper)

	mockIndexName := MockIndexNameExisting()

	type fields struct {
		config    HelperConfig
		logHelper log.IHelper
		rawClient *goES.Client
	}
	type args struct {
		index string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *DeleteIndexResp
		wantErr bool
	}{
		{
			fields: fields{
				config:    config,
				logHelper: logHelper,
				rawClient: rawClient,
			},
			args: args{
				index: mockIndexName,
			},
			want: &DeleteIndexResp{
				Acknowledged: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &helper{
				config:    tt.fields.config,
				logHelper: tt.fields.logHelper,
				rawClient: tt.fields.rawClient,
			}
			got, err := h.DeleteIndices(tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("helper.DeleteIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want.Acknowledged, got.Acknowledged)
		})
	}
}

func Test_helper_Reindex(t *testing.T) {
	config := MockHelperConfig()
	logHelper := tdLog.MockHelper()
	rawClient := MockRawClient(config, logHelper)

	mockSrcIndexName := MockIndexNameExisting()
	mockDestIndexName := MockIndexNameExisting()

	type fields struct {
		config    HelperConfig
		logHelper log.IHelper
		rawClient *goES.Client
	}
	type args struct {
		srcIndexName  string
		destIndexName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			fields: fields{
				config:    config,
				logHelper: logHelper,
				rawClient: rawClient,
			},
			args: args{
				srcIndexName:  mockSrcIndexName,
				destIndexName: mockDestIndexName,
			},
			want:    (*ReindexResp)(nil),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &helper{
				config:    tt.fields.config,
				logHelper: tt.fields.logHelper,
				rawClient: tt.fields.rawClient,
			}
			got, err := h.Reindex(tt.args.srcIndexName, tt.args.destIndexName)
			if (err != nil) != tt.wantErr {
				t.Errorf("helper.Reindex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("helper.Reindex() = %v, want %v", got, tt.want)
			}
		})
	}
}
