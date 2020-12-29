package log

import (
	"encoding/json"
	"esctl/pkg/config/dotenv"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type HelperConfig struct {
	LogLevel  string `env:"LOG_LEVEL,default=debug"`
	LogFormat string `env:"LOG_FORMAT,default=json"`
}

type IHelper interface {
	Debug(message string, fields map[string]interface{})
	Info(message string, fields map[string]interface{})
	Warn(message string, fields map[string]interface{})
	Error(message string, fields map[string]interface{})
	Fatal(message string, fields map[string]interface{})
	Panic(message string, fields map[string]interface{})
	LogLevel() string
	LogLevelNum() uint8
	SetWithField(key string, value interface{}) *helper
	SetWithGlobalField(key string, value interface{}) *helper
	NewChild() *helper
}

type helper struct {
	config           *HelperConfig
	logger           zerolog.Logger
	withFields       map[string]interface{}
	withGlobalFields map[string]interface{}
}

func DefaultHelperConfig() (*HelperConfig, error) {
	cfg := &HelperConfig{}
	if err := dotenv.Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func DefaultHelper() (*helper, error) {
	cfg, err := DefaultHelperConfig()
	if err != nil {
		return nil, errors.Wrap(err, "default config error")
	}
	return NewHelper(cfg)
}

func NewHelper(config *HelperConfig) (*helper, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	var logger zerolog.Logger

	// 设置日志格式
	if config.LogFormat != LOG_FORMAT_JSON {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		logger = zerolog.New(os.Stdout)
	}

	defaultWithFields := logger.With().Timestamp()

	// 设置是否输出文件路径、行号
	if config.LogLevel == LOG_LEVEL_ALL {
		// 由于封装的缘故，这个需要跳过前面 3 个帧
		defaultWithFields = defaultWithFields.CallerWithSkipFrameCount(3)
	}

	logger = defaultWithFields.Logger()

	switch config.LogLevel {
	case LOG_LEVEL_ALL:
		logger = logger.Level(zerolog.NoLevel)
	case LOG_LEVEL_DEBUG:
		logger = logger.Level(zerolog.DebugLevel)
	case LOG_LEVEL_INFO:
		logger = logger.Level(zerolog.InfoLevel)
	case LOG_LEVEL_WARN:
		logger = logger.Level(zerolog.WarnLevel)
	case LOG_LEVEL_ERROR:
		logger = logger.Level(zerolog.ErrorLevel)
	case LOG_LEVEL_FATAL:
		logger = logger.Level(zerolog.FatalLevel)
	case LOG_LEVEL_PANIC:
		logger = logger.Level(zerolog.PanicLevel)
	case LOG_LEVEL_DISABLED:
		logger = logger.Level(zerolog.Disabled)
	default:
		return nil, errors.Errorf("config.LogLevel=%s is invalid", config.LogLevel)
	}

	inst := &helper{
		config:           config,
		logger:           logger,
		withFields:       map[string]interface{}{},
		withGlobalFields: map[string]interface{}{},
	}

	return inst, nil
}

func (h *helper) Debug(message string, fields map[string]interface{}) {
	fields = h.mergeFileds(fields)
	h.logger.Debug().Fields(fields).Msg(message)
}

func (h *helper) Info(message string, fields map[string]interface{}) {
	fields = h.mergeFileds(fields)
	h.logger.Info().Fields(fields).Msg(message)
}

func (h *helper) Warn(message string, fields map[string]interface{}) {
	fields = h.mergeFileds(fields)
	h.logger.Warn().Fields(fields).Msg(message)
}

func (h *helper) Error(message string, fields map[string]interface{}) {
	fields = h.mergeFileds(fields)
	h.logger.Error().Fields(fields).Msg(message)
}

func (h *helper) Fatal(message string, fields map[string]interface{}) {
	fields = h.mergeFileds(fields)
	h.logger.Fatal().Fields(fields).Msg(message)
}

func (h *helper) Panic(message string, fields map[string]interface{}) {
	fields = h.mergeFileds(fields)
	h.logger.Panic().Fields(fields).Msg(message)
}

func (h *helper) LogLevel() string {
	return h.config.LogLevel
}

func (h *helper) LogLevelNum() uint8 {
	return LogLevelMapToNum[h.config.LogLevel]
}

func (h *helper) SetWithField(key string, value interface{}) *helper {
	if value != nil {
		h.withFields[key] = value
	} else {
		delete(h.withFields, key)
	}
	return h
}

func (h *helper) SetWithGlobalField(key string, value interface{}) *helper {
	if value != nil {
		h.withGlobalFields[key] = value
	} else {
		delete(h.withGlobalFields, key)
	}

	return h
}

func (h *helper) mergeFileds(customFields map[string]interface{}) map[string]interface{} {
	needMergedFields := []map[string]interface{}{}
	if h.withGlobalFields != nil {
		needMergedFields = append(needMergedFields, h.withGlobalFields)
	}
	if h.withFields != nil {
		needMergedFields = append(needMergedFields, h.withFields)
	}
	needMergedFields = append(needMergedFields, customFields)

	res := map[string]interface{}{}
	for _, fields := range needMergedFields {
		for k, v := range fields {
			res[k] = v
		}
	}

	return res
}

func (h *helper) NewChild() *helper {
	newInst := &helper{
		logger:           h.logger,
		config:           h.config,
		withFields:       map[string]interface{}{},
		withGlobalFields: h.withGlobalFields,
	}

	// 拷贝 map
	jsonStr, _ := json.Marshal(h.withFields)
	json.Unmarshal(jsonStr, &newInst.withFields)

	return newInst
}
