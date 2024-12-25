package middleware

import (
	"time"

	"github.com/BYT0723/go-tools/log"
	"github.com/BYT0723/go-tools/uctx"
	"github.com/gin-gonic/gin"
)

var (
	_ gin.HandlerFunc = WithTraceLogger(nil)
	_ gin.HandlerFunc = ApiLogger("")
)

func WithTraceLogger(logger log.Logger, fields ...log.Field) func(*gin.Context) {
	logger = logger.With(fields...)
	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(uctx.WithLogger(ctx, logger))
		ctx.Next()
	}
}

// must be after WithTraceLogger
// else do nothing
func ApiLogger(level string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var (
			l     = uctx.Logger(ctx)
			start = time.Now() // 开始时间
		)
		// 处理请求
		ctx.Next()

		// 结束时间
		latency := time.Since(start)

		l.With(
			log.Any("method", ctx.Request.Method),
			log.Any("path", ctx.Request.URL.Path),
			log.Any("status", ctx.Writer.Status()),
			log.Any("latency", latency),
		).Log(level, ctx.Request.Method)
	}
}
