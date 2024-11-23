package uctx

import (
	"context"

	"github.com/BYT0723/go-tools/log"
)

type uctxLoggerKey struct{}

// log.Logger or nil
func Logger(ctx context.Context) log.Logger {
	logger, _ := ctx.Value(uctxLoggerKey{}).(log.Logger)
	return logger
}

func WithLogger(ctx context.Context, logger log.Logger) context.Context {
	return context.WithValue(ctx, uctxLoggerKey{}, logger)
}
