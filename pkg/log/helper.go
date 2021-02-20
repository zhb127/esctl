package log

import (
	"encoding/json"

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
	config           HelperConfig
	rawLogger        *zerolog.Logger
	withFields       map[string]interface{}
	withGlobalFields map[string]interface{}
}

func NewHelper(config HelperConfig) (IHelper, error) {
	rawLogger, err := newRawLogger(config)
	if err != nil {
		return nil, err
	}

	inst := &helper{
		config:           config,
		rawLogger:        rawLogger,
		withFields:       map[string]interface{}{},
		withGlobalFields: map[string]interface{}{},
	}

	return inst, nil
}

func (h *helper) Debug(message string, fields map[string]interface{}) {
	fields = h.mergeFileds(fields)
	h.rawLogger.Debug().Fields(fields).Msg(message)
}

func (h *helper) Info(message string, fields map[string]interface{}) {
	fields = h.mergeFileds(fields)
	h.rawLogger.Info().Fields(fields).Msg(message)
}

func (h *helper) Warn(message string, fields map[string]interface{}) {
	fields = h.mergeFileds(fields)
	h.rawLogger.Warn().Fields(fields).Msg(message)
}

func (h *helper) Error(message string, fields map[string]interface{}) {
	fields = h.mergeFileds(fields)
	h.rawLogger.Error().Fields(fields).Msg(message)
}

func (h *helper) Fatal(message string, fields map[string]interface{}) {
	fields = h.mergeFileds(fields)
	h.rawLogger.Fatal().Fields(fields).Msg(message)
}

func (h *helper) Panic(message string, fields map[string]interface{}) {
	fields = h.mergeFileds(fields)
	h.rawLogger.Panic().Fields(fields).Msg(message)
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

	result := map[string]interface{}{}
	for _, fields := range needMergedFields {
		for k, v := range fields {
			result[k] = v
		}
	}

	return result
}

func (h *helper) NewChild() *helper {
	newRawLogger := *h.rawLogger

	newInst := &helper{
		config:           h.config,
		rawLogger:        &newRawLogger,
		withFields:       map[string]interface{}{},
		withGlobalFields: h.withGlobalFields,
	}

	// 拷贝 map
	jsonStr, _ := json.Marshal(h.withFields)
	_ = json.Unmarshal(jsonStr, &newInst.withFields)

	return newInst
}
