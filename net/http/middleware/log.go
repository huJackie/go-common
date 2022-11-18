package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/cloudadrd/go-common/code"
	"github.com/cloudadrd/go-common/log"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	const (
		format           = "%3d | %13v | %-15s | %-7s %s | %s %s | %d %s"
		defaultMaxMemory = 32 << 20 // 32 MB
	)
	return func(c *gin.Context) {
		var (
			start   = time.Now()
			ip      = c.ClientIP()
			path    = c.Request.URL.Path
			method  = c.Request.Method
			content = c.ContentType()
			params  string
			err     error
		)

		c.Error(code.ApiNotFound)

		switch {
		case strings.Contains(content, "multipart/form-data"):
			err = c.Request.ParseMultipartForm(defaultMaxMemory)
		default:
			err = c.Request.ParseForm()
		}
		if err != nil {
			c.Error(err)
		}

		switch method {
		case http.MethodPost:
			switch c.ContentType() {
			case "application/json":
				// 当post请求则忽略url中的参数
				byteData, err := ioutil.ReadAll(c.Request.Body)
				if err != nil {
					c.Error(err)
				}
				c.Request.Body = ioutil.NopCloser(bytes.NewReader(byteData))
				params = string(byteData)
			default:
				params = c.Request.Form.Encode()
			}
		case http.MethodGet:
			params = c.Request.Form.Encode()
		}

		c.Next()

		var (
			ctxErr  = c.Errors.ByType(gin.ErrorTypePrivate).Last()
			output  = log.Info
			latency = time.Since(start)
			isSlow  = latency >= time.Second
		)

		if ctxErr != nil {
			e, ok := code.Cause(ctxErr.Err).(code.Codes)
			if !ok {
				// todo:fixme 返回动态错误值?
				// e = code.NewError(code.ServeErr.Code(), ctxErr.Error())
				log.Warnf("Unknown err:%s", ctxErr.Err)
				e = code.ServeErr
			}
			if isSlow {
				output = log.Warn
			}
			output(fmt.Sprintf(
				format,
				c.Writer.Status(),
				latency,
				ip,
				method,
				path,
				params,
				c.GetString("uid"),
				e.Code(),
				e.Error(),
			))
		}
	}
}
