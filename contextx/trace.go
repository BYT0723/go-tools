package ctxx

import "context"

type traceID struct{}

func TraceID(ctx context.Context) (string, error) {
	key, ok := ctx.Value(traceID{}).(string)
	if !ok {
		return "", ErrApiKeyNotFound
	}
	return key, nil
}

func WithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, traceID{}, id)
}

type spanID struct{}

func SpanID(ctx context.Context) (string, error) {
	key, ok := ctx.Value(spanID{}).(string)
	if !ok {
		return "", ErrApiKeyNotFound
	}
	return key, nil
}

func WithSpanID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, spanID{}, id)
}

type requestID struct{}

func RequestID(ctx context.Context) (string, error) {
	key, ok := ctx.Value(requestID{}).(string)
	if !ok {
		return "", ErrApiKeyNotFound
	}
	return key, nil
}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestID{}, id)
}

type serviceName struct{}

func Service(ctx context.Context) (string, error) {
	key, ok := ctx.Value(serviceName{}).(string)
	if !ok {
		return "", ErrApiKeyNotFound
	}
	return key, nil
}

func WithService(ctx context.Context, service string) context.Context {
	return context.WithValue(ctx, serviceName{}, service)
}

type versionStr struct{}

func Version(ctx context.Context) (string, error) {
	key, ok := ctx.Value(versionStr{}).(string)
	if !ok {
		return "", ErrApiKeyNotFound
	}
	return key, nil
}

func WithVersion(ctx context.Context, version string) context.Context {
	return context.WithValue(ctx, versionStr{}, version)
}

type envStr struct{}

func Env(ctx context.Context) (string, error) {
	key, ok := ctx.Value(envStr{}).(string)
	if !ok {
		return "", ErrApiKeyNotFound
	}
	return key, nil
}

func WithEnv(ctx context.Context, env string) context.Context {
	return context.WithValue(ctx, envStr{}, env)
}
