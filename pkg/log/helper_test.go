package log

import (
	"esctl/pkg/config/dotenv"
	"reflect"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func MockHelperConfig() HelperConfig {
	config := HelperConfig{}
	if err := dotenv.Decode(&config); err != nil {
		panic(err)
	}
	return config
}

func TestNewHelper(t *testing.T) {
	mockHelperConfig := MockHelperConfig()

	type args struct {
		config HelperConfig
	}
	tests := []struct {
		name    string
		args    args
		want    IHelper
		wantErr bool
	}{
		{
			args: args{
				config: HelperConfig{
					LogLevel: "mock-error",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				config: mockHelperConfig,
			},
			want:    &helper{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHelper(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHelper() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.IsType(t, tt.want, got)
			if got != nil {
				assert.Implements(t, (*IHelper)(nil), got)
			}
		})
	}
}

func Test_helper_LogLevel(t *testing.T) {
	mockHelperConfig := MockHelperConfig()

	type fields struct {
		config           HelperConfig
		rawLogger        *zerolog.Logger
		withFields       map[string]interface{}
		withGlobalFields map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				config: mockHelperConfig,
			},
			want: mockHelperConfig.LogLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &helper{
				config:           tt.fields.config,
				rawLogger:        tt.fields.rawLogger,
				withFields:       tt.fields.withFields,
				withGlobalFields: tt.fields.withGlobalFields,
			}
			if got := h.LogLevel(); got != tt.want {
				t.Errorf("helper.LogLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_helper_LogLevelNum(t *testing.T) {
	mockHelperConfig := MockHelperConfig()

	type fields struct {
		config           HelperConfig
		rawLogger        *zerolog.Logger
		withFields       map[string]interface{}
		withGlobalFields map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   uint8
	}{
		{
			fields: fields{
				config: mockHelperConfig,
			},
			want: LogLevelMapToNum[mockHelperConfig.LogLevel],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &helper{
				config:           tt.fields.config,
				rawLogger:        tt.fields.rawLogger,
				withFields:       tt.fields.withFields,
				withGlobalFields: tt.fields.withGlobalFields,
			}
			if got := h.LogLevelNum(); got != tt.want {
				t.Errorf("helper.LogLevelNum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_helper_SetWithField(t *testing.T) {
	mockHelperConfig := MockHelperConfig()
	mockKey := "mock-key"
	mockValue := "mock-val"

	type fields struct {
		config           HelperConfig
		rawLogger        *zerolog.Logger
		withFields       map[string]interface{}
		withGlobalFields map[string]interface{}
	}
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]interface{}
	}{
		{
			fields: fields{
				config:     mockHelperConfig,
				withFields: map[string]interface{}{},
			},
			args: args{
				key:   mockKey,
				value: mockValue,
			},
			want: map[string]interface{}{
				mockKey: mockValue,
			},
		},
		{
			fields: fields{
				config: mockHelperConfig,
				withFields: map[string]interface{}{
					mockKey: mockValue,
				},
			},
			args: args{
				key:   mockKey,
				value: nil,
			},
			want: map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &helper{
				config:           tt.fields.config,
				rawLogger:        tt.fields.rawLogger,
				withFields:       tt.fields.withFields,
				withGlobalFields: tt.fields.withGlobalFields,
			}
			got := h.SetWithField(tt.args.key, tt.args.value)

			assert.Equal(t, tt.want, got.withFields)
		})
	}
}

func Test_helper_SetWithGlobalField(t *testing.T) {
	mockHelperConfig := MockHelperConfig()
	mockKey := "mock-key"
	mockValue := "mock-val"

	type fields struct {
		config           HelperConfig
		rawLogger        *zerolog.Logger
		withFields       map[string]interface{}
		withGlobalFields map[string]interface{}
	}
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]interface{}
	}{
		{
			fields: fields{
				config:           mockHelperConfig,
				withGlobalFields: map[string]interface{}{},
			},
			args: args{
				key:   mockKey,
				value: mockValue,
			},
			want: map[string]interface{}{
				mockKey: mockValue,
			},
		},
		{
			fields: fields{
				config: mockHelperConfig,
				withGlobalFields: map[string]interface{}{
					mockKey: mockValue,
				},
			},
			args: args{
				key:   mockKey,
				value: nil,
			},
			want: map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &helper{
				config:           tt.fields.config,
				rawLogger:        tt.fields.rawLogger,
				withFields:       tt.fields.withFields,
				withGlobalFields: tt.fields.withGlobalFields,
			}

			got := h.SetWithGlobalField(tt.args.key, tt.args.value)

			assert.Equal(t, tt.want, got.withGlobalFields)
		})
	}
}

func Test_helper_mergeFileds(t *testing.T) {
	mockHelperConfig := MockHelperConfig()
	mockKey := "mock-key"
	mockValue := "mock-val"
	mockGlobalKey := "mock-global-key"
	mockGlobalValue := "mock-global-val"
	mockCustomKey := "mock-custom-key"
	mockCustomValue := "mock-custom-value"

	type fields struct {
		config           HelperConfig
		rawLogger        *zerolog.Logger
		withFields       map[string]interface{}
		withGlobalFields map[string]interface{}
	}
	type args struct {
		customFields map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]interface{}
	}{
		{
			fields: fields{
				config: mockHelperConfig,
				withFields: map[string]interface{}{
					mockKey: mockValue,
				},
				withGlobalFields: map[string]interface{}{
					mockGlobalKey: mockGlobalValue,
				},
			},
			args: args{
				customFields: map[string]interface{}{
					mockCustomKey: mockCustomValue,
				},
			},
			want: map[string]interface{}{
				mockKey:       mockValue,
				mockGlobalKey: mockGlobalValue,
				mockCustomKey: mockCustomValue,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &helper{
				config:           tt.fields.config,
				rawLogger:        tt.fields.rawLogger,
				withFields:       tt.fields.withFields,
				withGlobalFields: tt.fields.withGlobalFields,
			}
			if got := h.mergeFileds(tt.args.customFields); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("helper.mergeFileds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_helper_NewChild(t *testing.T) {
	mockHelperConfig := MockHelperConfig()
	mockKey := "mock-key"
	mockValue := "mock-val"
	mockKey2 := "mock-key2"
	mockValue2 := "mock-val2"
	mockGlobalKey := "mock-global-key"
	mockGlobalValue := "mock-global-val"
	mockGlobalKey2 := "mock-global-key2"
	mockGlobalValue2 := "mock-global-value2"

	type fields struct {
		config           HelperConfig
		rawLogger        *zerolog.Logger
		withFields       map[string]interface{}
		withGlobalFields map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   *helper
	}{
		{
			fields: fields{
				config:    mockHelperConfig,
				rawLogger: &zerolog.Logger{},
				withFields: map[string]interface{}{
					mockKey: mockValue,
				},
				withGlobalFields: map[string]interface{}{
					mockGlobalKey: mockGlobalValue,
				},
			},
			want: &helper{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &helper{
				config:           tt.fields.config,
				rawLogger:        tt.fields.rawLogger,
				withFields:       tt.fields.withFields,
				withGlobalFields: tt.fields.withGlobalFields,
			}
			got := h.NewChild()
			assert.Equal(t, h.withGlobalFields, got.withGlobalFields)

			got.SetWithField(mockKey2, mockValue2)
			got.SetWithGlobalField(mockGlobalKey2, mockGlobalValue2)

			assert.NotEqual(t, h.withFields, got.withFields)
			assert.Equal(t, h.withGlobalFields, got.withGlobalFields)
		})
	}
}
