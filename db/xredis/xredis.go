package xredis

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	Url          string        `json:"url" yaml:"url"`                       // redis://<user>:<password>@<host>:<port>/<db_number>
	MaxRetries   int           `json:"max_retries" yaml:"max_retries"`       // 3
	PoolSize     int           `json:"pool_size" yaml:"pool_size"`           // cpu*10
	MinIdleConns int           `json:"min_idle_conns" yaml:"min_idle_conns"` //
	MaxConnAge   time.Duration `json:"max_conn_age" yaml:"max_conn_age"`     // 永久
	ReadTimeout  time.Duration `json:"read_timeout" yaml:"read_timeout"`     // 3 seconds
	DialTimeout  time.Duration `json:"dial_timeout" yaml:"dial_timeout"`     // 5 seconds
	PoolTimeout  time.Duration `json:"pool_timeout" yaml:"pool_timeout"`     // ReadTimeout+1seconds
	IdleTimeout  time.Duration `json:"idle_timeout" yaml:"idle_timeout"`     // 5 minutes
}

func (c *Config) Valid() error {
	if c == nil {
		return errors.New("redis config is nil")
	}
	if c.Url == "" {
		return errors.New("redis config url is empty")
	}
	return nil
}

func NewRedis(c *Config) *redis.Client {
	if err := c.Valid(); err != nil {
		log.Fatal(err)
	}
	opts, err := redis.ParseURL(c.Url)
	if err != nil {
		log.Fatal(err)
	}

	opts.MaxRetries = c.MaxRetries
	opts.PoolSize = c.PoolSize
	opts.MinIdleConns = c.MinIdleConns
	opts.MaxConnAge = c.MaxConnAge
	opts.ReadTimeout = c.ReadTimeout
	opts.DialTimeout = c.DialTimeout
	opts.PoolTimeout = c.PoolTimeout
	opts.IdleTimeout = c.IdleTimeout

	client := redis.NewClient(opts)
	if _, err := client.Ping(context.TODO()).Result(); err != nil {
		log.Fatalln(err)
	}
	return client
}
