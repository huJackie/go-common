package xtime

import "time"

func init() {
	local, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	time.Local = local
}
