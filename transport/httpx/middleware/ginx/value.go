package ginx

import (
	"context"

	ctxx "github.com/BYT0723/go-tools/contextx"
	"github.com/gin-gonic/gin"
)

var (
	_ gin.HandlerFunc = WithTraceID(func(ctx *gin.Context) string { return "" })
	_ gin.HandlerFunc = WithSpanID(func(ctx *gin.Context) string { return "" })
	_ gin.HandlerFunc = WithRequestID(func(ctx *gin.Context) string { return "" })
	_ gin.HandlerFunc = WithService("")
	_ gin.HandlerFunc = WithVersion("")
	_ gin.HandlerFunc = WithEnv("")
)

func WithValue(key any, generate func(ctx *gin.Context) any) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), key, generate(ctx)))
	}
}

func WithTraceID(idGenerate func(ctx *gin.Context) string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(
			ctxx.WithTraceID(ctx.Request.Context(), idGenerate(ctx)),
		)
		ctx.Next()
	}
}

func WithSpanID(idGenerate func(ctx *gin.Context) string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(
			ctxx.WithSpanID(ctx.Request.Context(), idGenerate(ctx)),
		)
		ctx.Next()
	}
}

func WithRequestID(idGenerate func(ctx *gin.Context) string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(
			ctxx.WithRequestID(ctx.Request.Context(), idGenerate(ctx)),
		)
		ctx.Next()
	}
}

func WithService(service string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(
			ctxx.WithService(ctx.Request.Context(), service),
		)
		ctx.Next()
	}
}

func WithVersion(version string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(
			ctxx.WithVersion(ctx.Request.Context(), version),
		)
		ctx.Next()
	}
}

func WithEnv(env string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(
			ctxx.WithEnv(ctx.Request.Context(), env),
		)
		ctx.Next()
	}
}
