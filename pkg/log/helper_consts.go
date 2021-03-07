package log

// 日志级别越小（从小到大：all -> disabled），输出的日志越多、越细
const (
	LogLevelAllNum = iota + 1
	LogLevelDebugNum
	LogLevelInfoNum
	LogLevelWarnNum
	LogLevelErrorNum
	LogLevelFatalNum
	LogLevelPanicNum
	LogLevelNoneNum
)

const (
	LogLevelAll   = "all"
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
	LogLevelFatal = "fatal"
	LogLevelPanic = "panic"
	LogLevelNone  = "none"
)

var LogLevelMapToNum = map[string]uint8{
	LogLevelAll:   LogLevelAllNum,
	LogLevelDebug: LogLevelDebugNum,
	LogLevelInfo:  LogLevelInfoNum,
	LogLevelWarn:  LogLevelWarnNum,
	LogLevelError: LogLevelErrorNum,
	LogLevelFatal: LogLevelFatalNum,
	LogLevelPanic: LogLevelPanicNum,
	LogLevelNone:  LogLevelNoneNum,
}

const (
	LogFormatJSON   = "json"
	LogFormatPretty = "pretty"
)
