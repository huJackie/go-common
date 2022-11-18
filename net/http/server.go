package http

import (
	"errors"
	"log"
	"net"
	"net/http"
	"time"
)

var DefaultConfig = &Config{
	Network:      "tcp",
	Addr:         ":8080",
	ReadTimeout:  time.Second * 5,
	WriteTimeout: time.Second * 10,
	IdleTimeout:  time.Second * 10,
	Monitor:      true,
}

type Config struct {
	Network      string        `json:"network" yaml:"network"`             // tcp ..
	Addr         string        `json:"addr" yaml:"addr"`                   // http://127.0.0.1:8080
	ReadTimeout  time.Duration `json:"read_timeout" yaml:"read_timeout"`   // 读取超时
	WriteTimeout time.Duration `json:"write_timeout" yaml:"write_timeout"` // 写超时
	IdleTimeout  time.Duration `json:"idle_timeout" yaml:"idle_timeout"`   // 空闲连接超时
	Monitor      bool          `json:"monitor" yaml:"monitor"`             // 是否开启日志 pprof trace log...
}

func (c *Config) Valid() error {
	if c == nil {
		return errors.New("http.config is nil")
	}
	return nil
}

func Run(handler http.Handler, conf *Config) error {
	traces(handler, conf.Monitor)

	server := &http.Server{
		Handler:      handler,
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
		IdleTimeout:  conf.IdleTimeout,
	}
	//server.SetKeepAlivesEnabled(true)

	l, err := net.Listen(conf.Network, conf.Addr)
	if err != nil {
		return err
	}

	go func() {
		if err := server.Serve(l); err != nil {
			if err == http.ErrServerClosed {
				log.Printf("xhttp [%s] server exit...\n", conf.Addr)
				return
			}
			log.Panicf("xhttp Serve err:%s\n", err)
		}
	}()
	log.Printf("xhttp server running on [%s]\n", conf.Addr)
	return nil
}
