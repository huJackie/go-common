package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
	设置日志级别
	case条件越往下日志级别越高越严重日志等级如下:
	debug < info < warn < error < dPanic < panic < fatal
	zap默认级别是info级别,当我们调整日志等级比如调整到warn级别的日志,此时debug info则不会输出
*/
func atomicLevel(l string) zap.AtomicLevel {
	var level zapcore.Level
	switch l {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
	return zap.NewAtomicLevelAt(level)
}