package ginmd

import (
	uctx "github.com/BYT0723/go-tools/ctx"
	"github.com/BYT0723/go-tools/log"
	"github.com/gin-gonic/gin"
)

func WithTraceLogger(logger log.Logger, fields ...*log.Field) func(*gin.Context) {
	logger = logger.With(fields...)
	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(uctx.WithLogger(ctx, logger))
		ctx.Next()
	}
}
