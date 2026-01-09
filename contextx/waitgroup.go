package ctxx

import (
	"context"
	"sync"
)

type srvxWgKey struct{}

func WithWaitGroup(ctx context.Context, wg *sync.WaitGroup) context.Context {
	return context.WithValue(ctx, srvxWgKey{}, wg)
}

func WaitGroup(ctx context.Context) (*sync.WaitGroup, error) {
	wg, ok := ctx.Value(srvxWgKey{}).(*sync.WaitGroup)
	if !ok {
		return nil, ErrWaitGroupNotFound
	}
	return wg, nil
}
