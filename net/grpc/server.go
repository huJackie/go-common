package grpc

import (
	"time"
)

var DefaultServerConfig = &ServerConfig{
	Name:              "default",
	Network:           "tcp",
	Addr:              ":50051",
	Trace:             false,
	Timeout:           time.Second,
	IdleTimeout:       0,
	MaxLifeTime:       0,
	ForceCloseWait:    0,
	KeepAliveInterval: 0,
	KeepAliveTimeout:  0,
}

type ServerConfig struct {
	Name              string        // 服务名称
	Network           string        // tcp ...
	Addr              string        // 地址
	Trace             bool          // 是否开启trace
	Timeout           time.Duration // 超时(SetDeadline())
	IdleTimeout       time.Duration // xtime.Duration(time.Second * 60),
	MaxLifeTime       time.Duration // xtime.Duration(time.Hour * 2),
	ForceCloseWait    time.Duration // xtime.Duration(time.Second * 20),
	KeepAliveInterval time.Duration // xtime.Duration(time.Second * 60),
	KeepAliveTimeout  time.Duration // xtime.Duration(time.Second * 20),
}
