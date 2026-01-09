package ctxx

import (
	"context"
	"net"
)

type listenerKey struct{}

func WithListener(ctx context.Context, l net.Listener) context.Context {
	return context.WithValue(ctx, listenerKey{}, l)
}

func Listener(ctx context.Context) (net.Listener, error) {
	l, ok := ctx.Value(listenerKey{}).(net.Listener)
	if !ok {
		return nil, ErrListenerNotFound
	}
	return l, nil
}
