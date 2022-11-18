package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"hash/fnv"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cloudadrd/go-common/code"
	"github.com/cloudadrd/go-common/log"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

func Json(ctx *gin.Context, data interface{}, err error) {
	if err == nil {
		err = code.OK
	}
	ctx.Error(err)
	e, ok := code.Cause(err).(code.Codes)
	if !ok {
		log.Debugf("json err: %s", err)
		result := map[string]interface{}{
			"code":   code.ServeErr.Code(),
			"msg":    err.Error(),
			"data":   data,
			"detail": err.Error(),
		}
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := map[string]interface{}{
		"code":   e.Code(),
		"msg":    e.Error(),
		"data":   data,
		"detail": err.Error(),
	}
	ctx.JSON(http.StatusOK, result)
}

type Config interface {
	Valid() error
}

func Valid(c ...Config) error {
	for _, f := range c {
		if err := f.Valid(); err != nil {
			return err
		}
	}
	return nil
}

//传入字符串百分比转换成对应的小数
func Percent(per string) float64 {
	temp := strings.ReplaceAll(per, "%", "")
	result, _ := strconv.ParseFloat(temp, 10)
	return result
}

// 读取配置文件
func Unmarshal(path string, out interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, out); err != nil {
		return err
	}
	if o, ok := out.(interface{ Valid() error }); ok {
		return o.Valid()
	}
	return errors.New("invalid type")
}

func ExpireTime(day int, t ...time.Time) time.Duration {
	var (
		h int
		m int
		s int
	)
	if t != nil && len(t) > 0 {
		h, m, s = t[0].Clock()
	} else {
		h, m, s = time.Now().Clock()
	}
	return (time.Duration((day * 86400) + 86400 - (h*3600 + (m * 60) + s))) * time.Second
}

func Anyone(array []string, target string) bool {
	for _, v := range array {
		if strings.ToLower(v) == strings.ToLower(target) {
			return true
		}
	}
	return false
}

func AnyoneInt(array []int64, target int64) bool {
	for _, v := range array {
		if v == target {
			return true
		}
	}
	return false
}

func Now() time.Time {
	l, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(l)
}

func HashCode(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// 慎用
func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
