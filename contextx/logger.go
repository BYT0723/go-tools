package ctxx

import (
	"context"

	"github.com/BYT0723/go-tools/logx"
)

type loggerKey struct{}

// log.Logger or nil
func Logger(ctx context.Context) (logx.Logger, error) {
	l, ok := ctx.Value(loggerKey{}).(logx.Logger)
	if !ok {
		return nil, ErrLoggerNotFound
	}
	return l, nil
}

func WithLogger(ctx context.Context, logger logx.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}
