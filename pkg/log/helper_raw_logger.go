package log

import (
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func newRawLogger(config HelperConfig) (*zerolog.Logger, error) {
	var logger zerolog.Logger

	// 设置日志输出格式
	if config.LogFormat != LogFormatJSON {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		logger = zerolog.New(os.Stdout)
	}

	withDefaultFields := logger.With().Timestamp()

	// 设置是否输出文件路径、行号
	if config.LogLevel == LogLevelAll {
		// 由于封装的缘故，这个需要跳过前面 3 个帧
		withDefaultFields = withDefaultFields.CallerWithSkipFrameCount(3)
	}

	logger = withDefaultFields.Logger()

	switch config.LogLevel {
	case LogLevelAll:
		logger = logger.Level(zerolog.NoLevel)
	case LogLevelDebug:
		logger = logger.Level(zerolog.DebugLevel)
	case LogLevelInfo:
		logger = logger.Level(zerolog.InfoLevel)
	case LogLevelWarn:
		logger = logger.Level(zerolog.WarnLevel)
	case LogLevelError:
		logger = logger.Level(zerolog.ErrorLevel)
	case LogLevelFatal:
		logger = logger.Level(zerolog.FatalLevel)
	case LogLevelPanic:
		logger = logger.Level(zerolog.PanicLevel)
	case LogLevelNone:
		logger = logger.Level(zerolog.Disabled)
	default:
		return nil, errors.Errorf("config.LogLevel=%s is invalid", config.LogLevel)
	}

	return &logger, nil
}
