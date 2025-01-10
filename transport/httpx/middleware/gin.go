package middleware

import (
	"time"

	"github.com/BYT0723/go-tools/contextx"
	"github.com/BYT0723/go-tools/logx"
	"github.com/gin-gonic/gin"
)

var (
	_ gin.HandlerFunc = WithTraceLogger(nil)
	_ gin.HandlerFunc = ApiLogger("")
	_ gin.HandlerFunc = WithApiKey(func(ctx *gin.Context) string { return "" })
)

func WithTraceLogger(logger logx.Logger) func(*gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(contextx.WithLogger(ctx, logger))
		ctx.Next()
	}
}

func WithApiKey(keyGenerate func(ctx *gin.Context) string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(contextx.WithApiKey(ctx, keyGenerate(ctx)))
		ctx.Next()
	}
}

// must be after WithTraceLogger
// else do nothing
func ApiLogger(level string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var (
			l     = contextx.Logger(ctx)
			start = time.Now() // 开始时间
		)
		// 处理请求
		ctx.Next()

		// 结束时间
		latency := time.Since(start)

		l.With(
			logx.Any("method", ctx.Request.Method),
			logx.Any("path", ctx.Request.URL.Path),
			logx.Any("status", ctx.Writer.Status()),
			logx.Any("latency", latency),
		).Log(level, ctx.Request.Method)
	}
}
