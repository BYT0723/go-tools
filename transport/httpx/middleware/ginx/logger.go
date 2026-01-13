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
)

// must be after WithTraceLogger
// else do nothing
func WithApiLog(level string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		// 开始时间
		start := time.Now()

		// 处理请求
		ctx.Next()

		l, err := ctxx.Logger(ctx.Request.Context())
		if err != nil {
			return
		}

		var (
			// 字段
			fields = make([]logx.Field, 0, 10)
			// 结束时间
			latency = time.Since(start)
		)

		if traceID, err := ctxx.TraceID(ctx.Request.Context()); err == nil {
			fields = append(fields, logx.String("trace_id", traceID))
		}
		if spanID, err := ctxx.SpanID(ctx.Request.Context()); err == nil {
			fields = append(fields, logx.String("span_id", spanID))
		}
		if reqID, err := ctxx.RequestID(ctx.Request.Context()); err == nil {
			fields = append(fields, logx.String("request_id", reqID))
		}
		if service, err := ctxx.Service(ctx.Request.Context()); err == nil {
			fields = append(fields, logx.String("service", service))
		}
		if version, err := ctxx.Version(ctx.Request.Context()); err == nil {
			fields = append(fields, logx.String("version", version))
		}
		if env, err := ctxx.Env(ctx.Request.Context()); err == nil {
			fields = append(fields, logx.String("env", env))
		}

		fields = append(fields,
			logx.String("method", ctx.Request.Method),
			logx.String("path", ctx.Request.URL.Path),
			logx.Int("status", ctx.Writer.Status()),
			logx.Duration("latency", latency),
		)
		l.Log(level, "API REQUEST", fields...)
	}
}

func WithTraceLogger(logger logx.Logger) func(*gin.Context) {
	return func(ctx *gin.Context) {
		if logger != nil {
			ctx.Request = ctx.Request.WithContext(ctxx.WithLogger(ctx.Request.Context(), logger))
		}
		ctx.Next()
	}
}
