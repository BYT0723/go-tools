package srvx

import (
	"context"
	"errors"
	"sync"
)

var ErrWaitGroupNotFound = errors.New("waitgroup not found")

type srvxWgKey struct{}

func WithWaitGroup(ctx context.Context, wg *sync.WaitGroup) context.Context {
	return context.WithValue(ctx, srvxWgKey{}, wg)
}

func WaitGroup(ctx context.Context) (*sync.WaitGroup, error) {
	v := ctx.Value(srvxWgKey{})
	if v == nil {
		return nil, ErrWaitGroupNotFound
	}
	wg, ok := v.(*sync.WaitGroup)
	if !ok {
		return nil, ErrWaitGroupNotFound
	}
	return wg, nil
}
