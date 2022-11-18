package env

import (
	"flag"
	"os"
)

var (
	Hostname string // hostname
	Mode     string // run mode
)

const (
	DEBUG   = "debug"   // 开发
	TEST    = "test"    // 测试
	GRAY    = "gray"    // 灰度
	RELEASE = "release" // 生产
)

func init() {
	var err error
	if Hostname, err = os.Hostname(); err != nil || Hostname == "" {
		Hostname = os.Getenv("HOSTNAME")
	}
	addFlag(flag.CommandLine)
}

func defaultString(env, value string) string {
	v := os.Getenv(env)
	if v == "" {
		return value
	}
	return v
}

// 优先级从低到高 default < goenv < command
func addFlag(f *flag.FlagSet) {
	f.StringVar(&Mode, "mode", defaultString("MODE", DEBUG), "run mode,default debug mode.value:[debug,test,gray,release]")
}
