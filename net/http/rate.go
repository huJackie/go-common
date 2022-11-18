package http

import (
	"time"

	"github.com/cloudadrd/go-common"
	"github.com/cloudadrd/go-common/code"
	"github.com/cloudadrd/go-common/rate/limit"

	"github.com/gin-gonic/gin"
)

/*
	限流
*/
func Limit(bust int32) gin.HandlerFunc {
	l := limit.NewLimiter(time.Second*1, bust)
	return func(context *gin.Context) {
		if !l.Increase() {
			pkg.Json(context, nil, code.Retry)
			context.Abort()
			return
		}
		context.Next()
	}
}
