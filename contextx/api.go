package ctxx

import "context"

type apiKey struct{}

func ApiKey(ctx context.Context) (string, error) {
	key, ok := ctx.Value(apiKey{}).(string)
	if !ok {
		return "", ErrApiKeyNotFound
	}
	return key, nil
}

func WithApiKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, apiKey{}, key)
}
