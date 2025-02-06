package contextx

import (
	"context"

	"github.com/BYT0723/go-tools/logx"
)

type contextxLoggerKey struct{}

// log.Logger or nil
func Logger(ctx context.Context) logx.Logger {
	logger, _ := ctx.Value(contextxLoggerKey{}).(logx.Logger)
	return logger
}

func WithLogger(ctx context.Context, logger logx.Logger) context.Context {
	return context.WithValue(ctx, contextxLoggerKey{}, logger)
}
