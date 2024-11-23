package middleware

import (
	"github.com/BYT0723/go-tools/log"
	"github.com/BYT0723/go-tools/uctx"
	"github.com/gin-gonic/gin"
)

var _ gin.HandlerFunc = WithTraceLogger(nil)

func WithTraceLogger(logger log.Logger, fields ...*log.Field) func(*gin.Context) {
	logger = logger.With(fields...)
	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(uctx.WithLogger(ctx, logger))
		ctx.Next()
	}
}
