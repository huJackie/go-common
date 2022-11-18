package log

import (
	"fmt"
	"path/filepath"

	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type level string

const (
	_info  level = "info"
	_warn        = "warn"
	_error       = "error"
)

// 日志分割
func split(c *Config, l level) zapcore.WriteSyncer {
	var (
		format = fmt.Sprintf("%s.log", l)
		name   = filepath.Join(c.Path, format)
	)
	rotate := &lumberjack.Logger{
		Filename:   name,         // 日志文件名字
		MaxSize:    c.MaxSize,    // 日志文件的大小 单位是M 默认是100M
		MaxAge:     c.MaxAge,     // 日志文件的生命周期
		MaxBackups: c.MaxBackups, // 保留旧的日志文件最大数量
		LocalTime:  true,         // 是否使用本地时间进行日志格式化等,默认false代表使用UTC时间
		Compress:   c.Compress,   // 开启gzip压缩 默认不开启
	}
	return zapcore.AddSync(rotate)
}
