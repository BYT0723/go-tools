package uctx

import (
	"context"
)

type apiKey struct{}

// log.Logger or nil
func ApiKey(ctx context.Context) string {
	key, _ := ctx.Value(apiKey{}).(string)
	return key
}

func WithApiKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, apiKey{}, key)
}
