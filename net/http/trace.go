package http

import (
	"expvar"
	"net/http"
	"sync"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/trace"
)

var once sync.Once

func handle(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func traces(handler http.Handler, onOff bool) {
	if !onOff {
		return
	}
	engine, ok := handler.(*gin.Engine)
	if !ok {
		panic("http.Handler is not *gin.Engine type")
	}
	once.Do(func() {
		//pprof
		pprof.Register(engine)
		//expvar
		engine.GET("/debug/vars", handle(expvar.Handler().ServeHTTP))
		//rpc
		engine.GET("/debug/requests", handle(trace.Traces))
		engine.GET("/debug/events", handle(trace.Events))
	})
}
