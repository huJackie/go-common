package xxorm

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

// Config mysql配置
type Config struct {
	User      string        `yaml:"user"`
	Password  string        `yaml:"password"`
	Protol    string        `yaml:"protol"`
	Host      string        `yaml:"host"`
	Port      int           `yaml:"port"`
	DbName    string        `yaml:"dbname"`
	MaxIdle   int           `yaml:"maxIdle"`
	MaxActive int           `yaml:"maxActive"`
	MaxLife   time.Duration `yaml:"maxLife"`
	ShowSQL   bool          `yaml:"showSQL"`
	Level     int           `yaml:"level"` // 0:debug 1:info 2:warn 3:error
}

func (c *Config) Valid() error {
	if c == nil {
		return fmt.Errorf("Config is nil")
	}
	return nil
}

// New
func New(c *Config) (*xorm.Engine, error) {
	if err := c.Valid(); err != nil {
		return nil, err
	}
	const url = "%s:%s@%s(%s:%d)/%s?charset=utf8mb4&parseTime=true"
	var addr = fmt.Sprintf(url, c.User, c.Password, c.Protol, c.Host, c.Port, c.DbName)
	engine, err := xorm.NewEngine("mysql", addr)
	if err != nil {
		return nil, fmt.Errorf("NewEngine():%w", err)
	}
	if err := engine.Ping(); err != nil {
		return nil, fmt.Errorf("Ping():%w", err)
	}

	engine.SetMaxIdleConns(c.MaxIdle)
	engine.SetMaxOpenConns(c.MaxActive)
	engine.ShowSQL(c.ShowSQL)
	engine.SetConnMaxLifetime(c.MaxLife)
	engine.SetLogLevel(log.LogLevel(c.Level))
	return engine, nil
}
