package log

import (
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func newRawLogger(config HelperConfig) (*zerolog.Logger, error) {
	var logger zerolog.Logger

	// 设置日志输出格式
	if config.LogFormat != LOG_FORMAT_JSON {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		logger = zerolog.New(os.Stdout)
	}

	withDefaultFields := logger.With().Timestamp()

	// 设置是否输出文件路径、行号
	if config.LogLevel == LOG_LEVEL_ALL {
		// 由于封装的缘故，这个需要跳过前面 3 个帧
		withDefaultFields = withDefaultFields.CallerWithSkipFrameCount(3)
	}

	logger = withDefaultFields.Logger()

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

	return &logger, nil
}
