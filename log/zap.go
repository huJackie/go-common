package log

import (
	"os"
	"time"

	"github.com/cloudadrd/go-common/conf/env"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger  *zap.SugaredLogger
	_logger *zap.Logger
	_sugar  *zap.SugaredLogger
	c       *Config
)

func init() {
	host, _ := os.Hostname()
	c = &Config{
		Level:      "info",
		Host:       host,
		AppId:      "1",
		Path:       "./log/",
		Stdout:     true,
		MaxSize:    100,
		MaxBackups: 0,
		MaxAge:     0,
		Compress:   false,
	}
}

func defaultString(old, new string) string {
	if new != "" {
		return new
	}
	return old
}

func New(conf *Config) {
	if conf != nil {
		c.Host = defaultString(c.Host, conf.Host)
		c.Level = defaultString(c.Host, atomicLevel(conf.Level).String())
		c.Path = defaultString(c.Path, conf.Path)
		c.AppId = defaultString(c.AppId, conf.AppId)
		c.Stdout = conf.Stdout
		c.MaxBackups = conf.MaxBackups
		c.MaxAge = conf.MaxAge
		c.MaxSize = conf.MaxSize
		c.Compress = conf.Compress
	}
	var (
		encoderConfig = zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			TimeKey:        "ts",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "stack",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.MillisDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		}
		lowLevel = zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= atomicLevel(c.Level).Level() && l < zapcore.ErrorLevel
		})
		highLevel = zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= zapcore.ErrorLevel
		})
		infoWrite   = []zapcore.WriteSyncer{split(c, _info)}
		errorWrite  = []zapcore.WriteSyncer{split(c, _error)}
		lowEncoder  zapcore.Encoder
		highEncoder zapcore.Encoder
	)

	if c.Stdout {
		infoWrite = append(infoWrite, zapcore.AddSync(os.Stdout))
		errorWrite = append(errorWrite, zapcore.AddSync(os.Stderr))
	}

	switch env.Mode {
	case env.RELEASE:
		lowEncoder = zapcore.NewJSONEncoder(encoderConfig)
		highEncoder = zapcore.NewJSONEncoder(encoderConfig)
	case env.DEBUG:
		fallthrough
	default:
		lowEncoder = zapcore.NewConsoleEncoder(encoderConfig)
		highEncoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core := zapcore.NewTee(
		zapcore.NewCore(lowEncoder, zapcore.NewMultiWriteSyncer(infoWrite...), lowLevel),
		zapcore.NewCore(highEncoder, zapcore.NewMultiWriteSyncer(errorWrite...), highLevel),
	)

	var opts = []zap.Option{
		zap.WithCaller(true),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(highLevel),
		zap.Fields(zap.String("appId", c.AppId), zap.String("host", c.Host)),
	}

	_logger = zap.New(zapcore.NewSamplerWithOptions(core, time.Second, 100, 100), opts...)
	_sugar = _logger.Sugar()
	Logger = _logger.Sugar()
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debugw(msg string, fields ...zap.Field) {
	_logger.Debug(msg, fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Infow(msg string, fields ...zap.Field) {
	_logger.Info(msg, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Warnw(msg string, fields ...zap.Field) {
	_logger.Warn(msg, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Errorw(msg string, fields ...zap.Field) {
	_logger.Error(msg, fields...)
}

// DPanic logs a message at DPanicLevel. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func DPanicw(msg string, fields ...zap.Field) {
	_logger.DPanic(msg, fields...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panicw(msg string, fields ...zap.Field) {
	_logger.Panic(msg, fields...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatalw(msg string, fields ...zap.Field) {
	_logger.Fatal(msg, fields...)
}

// Debug uses fmt.Sprint to construct and log a message.
func Debug(args ...interface{}) {
	_sugar.Debug(args...)
}

// Info uses fmt.Sprint to construct and log a message.
func Info(args ...interface{}) {
	_sugar.Info(args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(args ...interface{}) {
	_sugar.Warn(args...)
}

// Error uses fmt.Sprint to construct and log a message.
func Error(args ...interface{}) {
	_sugar.Error(args...)
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the
// logger then panics. (See DPanicLevel for details.)
func DPanic(args ...interface{}) {
	_sugar.DPanic(args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func Panic(args ...interface{}) {
	_sugar.Panic(args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	_sugar.Fatal(args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	_sugar.Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
	_sugar.Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	_sugar.Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	_sugar.Errorf(template, args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics. (See DPanicLevel for details.)
func DPanicf(template string, args ...interface{}) {
	_sugar.DPanicf(template, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(template string, args ...interface{}) {
	_sugar.Panicf(template, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(template string, args ...interface{}) {
	_sugar.Fatalf(template, args...)
}

//// Debugw logs a message with some additional context. The variadic key-value
//// pairs are treated as they are in With.
////
//// When debug-atomicLevel logging is disabled, this is much faster than
////  s.With(keysAndValues).Debug(msg)
//func Debugw(msg string, keysAndValues ...interface{}) {
//	_sugar.Debugw(msg, keysAndValues...)
//}
//
//// Infow logs a message with some additional context. The variadic key-value
//// pairs are treated as they are in With.
//func Infow(msg string, keysAndValues ...interface{}) {
//	_sugar.Infow(msg, keysAndValues...)
//}
//
//// Warnw logs a message with some additional context. The variadic key-value
//// pairs are treated as they are in With.
//func Warnw(msg string, keysAndValues ...interface{}) {
//	_sugar.Warnw(msg, keysAndValues...)
//}
//
//// Errorw logs a message with some additional context. The variadic key-value
//// pairs are treated as they are in With.
//func Errorw(msg string, keysAndValues ...interface{}) {
//	_sugar.Warnw(msg, keysAndValues...)
//}
//
//// DPanicw logs a message with some additional context. In development, the
//// logger then panics. (See DPanicLevel for details.) The variadic key-value
//// pairs are treated as they are in With.
//func DPanicw(msg string, keysAndValues ...interface{}) {
//	_sugar.DPanicw(msg, keysAndValues...)
//}
//
//// Panicw logs a message with some additional context, then panics. The
//// variadic key-value pairs are treated as they are in With.
//func Panicw(msg string, keysAndValues ...interface{}) {
//	_sugar.DPanicw(msg, keysAndValues...)
//}
//
//// Fatalw logs a message with some additional context, then calls os.Exit. The
//// variadic key-value pairs are treated as they are in With.
//func Fatalw(msg string, keysAndValues ...interface{}) {
//	_sugar.Fatalw(msg, keysAndValues...)
//}
