package ginx

import (
	"time"

	ctxx "github.com/BYT0723/go-tools/contextx"
	"github.com/BYT0723/go-tools/logx"
	"github.com/gin-gonic/gin"
)

var (
	_ gin.HandlerFunc = WithApiLog("info")
	_ gin.HandlerFunc = WithTraceLogger(nil)
	_ gin.HandlerFunc = WithApiKey(func(ctx *gin.Context) string { return "" })
)

// must be after WithTraceLogger
// else do nothing
func WithApiLog(level string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var (
			l, err    = ctxx.Logger(ctx.Request.Context())
			start     = time.Now() // 开始时间
			apiKey, _ = ctxx.ApiKey(ctx.Request.Context())
		)

		// 处理请求
		ctx.Next()

		if err != nil {
			return
		}

		// 结束时间
		latency := time.Since(start)

		l.With(
			logx.String("method", ctx.Request.Method),
			logx.String("path", ctx.Request.URL.Path),
			logx.Int("status", ctx.Writer.Status()),
			logx.String("api_key", apiKey),
			logx.Duration("latency", latency),
		).Log(level, "API Request")
	}
}

func WithTraceLogger(logger logx.Logger) func(*gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(ctxx.WithLogger(ctx.Request.Context(), logger))
		ctx.Next()
	}
}

func WithApiKey(keyGenerate func(ctx *gin.Context) string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(
			ctxx.WithApiKey(ctx.Request.Context(), keyGenerate(ctx)),
		)
		ctx.Next()
	}
}
