package jwt

import (
	"fmt"
	"time"
)

type Config struct {
	PrivateKey string        `json:"privateKey" yaml:"privateKey"`
	Expire     time.Duration `json:"expire" yaml:"expire"`
}

func (c *Config) Valid() error {
	if c == nil {
		return fmt.Errorf("Config is nil.")
	}
	if c.PrivateKey == "" {
		return fmt.Errorf("Config.PrivateKey is empty.")
	}
	return nil
}
