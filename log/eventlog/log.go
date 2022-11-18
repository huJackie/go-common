package eventlog

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Config struct {
	Stdout   bool   `json:"stdout" yaml:"stdout"`
	FileName string `json:"stdout" yaml:"filename"`
	Location string `json:"location" yaml:"location"`
	location *time.Location
}

func (c *Config) Valid() error {
	if c == nil {
		return errors.New("Config is nil. ")
	}
	c.location = time.Local
	if c.Location != "" {
		loc, err := time.LoadLocation(c.Location)
		if err != nil {
			return fmt.Errorf("Config.Location is valid. ")
		}
		c.location = loc
	}
	if filepath.Ext(c.FileName) == "" {
		return fmt.Errorf("Config.FileName is deficiency file ext. ")
	}
	return nil
}

type EventLog interface {
	Close(ctx context.Context) error
	Printf(string, ...interface{})
	Println(...interface{})
}

type logger struct {
	conf   *Config
	mu     sync.Mutex
	file   *os.File
	log    *log.Logger
	stdout *log.Logger
}

func New(c *Config) EventLog {
	if err := c.Valid(); err != nil {
		panic(fmt.Sprintf("log:%s", err))
	}

	l := &logger{
		conf: c,
	}
	if c.Stdout {
		l.stdout = log.New(os.Stdout, "[event] ", log.LstdFlags|log.Lshortfile)
	}
	if err := l.openFile(); err != nil {
		panic(fmt.Sprintf("openFile:%s", err))
	}
	l.log = log.New(l.file, "[event] ", log.LstdFlags|log.Lshortfile)
	go func() {
		if err := l.rotate(); err != nil {
			panic(fmt.Sprintf("rotate:%s", err))
		}
	}()
	return l
}

func (l *logger) Close(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.close()
}

// close closes the file if it is open.
func (l *logger) close() error {
	if l.file == nil {
		return nil
	}
	err := l.file.Close()
	l.file = nil
	return err
}

func (l *logger) openFile() error {
	if err := os.MkdirAll(filepath.Dir(l.conf.FileName), 0744); err != nil {
		return fmt.Errorf("can't make directories for new logfile: %w", err)
	}

	f, err := os.OpenFile(l.filename(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("can't open new logfile: %w", err)
	}
	l.file = f
	if l.log != nil {
		l.log.SetOutput(f)
	}
	return nil
}

func (l *logger) filename() string {
	day := time.Now().In(l.conf.location).Format("20060102")
	dir, name := filepath.Split(l.conf.FileName)
	ext := filepath.Ext(name)
	prefix := name[:len(name)-len(ext)]
	return filepath.Join(dir, fmt.Sprintf("%s-%s%s", prefix, day, ext))
}

func (l *logger) rotate() error {
	h, m, s := time.Now().In(l.conf.location).Clock()
	dur := time.Duration(86400-(h*3600+(m*60)+s)) * time.Second
	reset := time.Hour * 24

	tick := time.NewTimer(dur)
	defer tick.Stop()
	for range tick.C {
		tick.Reset(reset)
		handler := func() error {
			l.mu.Lock()
			defer l.mu.Unlock()
			if err := l.close(); err != nil {
				return fmt.Errorf("close():%w", err)
			}
			if err := l.openFile(); err != nil {
				return fmt.Errorf("openFile():%w", err)
			}
			return nil
		}
		if err := handler(); err != nil {
			return err
		}
	}
	return nil
}

func (l *logger) Printf(format string, args ...interface{}) {
	if l.stdout != nil {
		l.stdout.Output(2, fmt.Sprintf(format, args...))
	}
	l.log.Output(2, fmt.Sprintf(format, args...))
}

func (l *logger) Println(str ...interface{}) {
	if l.stdout != nil {
		l.stdout.Output(2, fmt.Sprintln(str...))
	}
	l.log.Output(2, fmt.Sprintln(str...))
}
