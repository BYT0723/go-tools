package channelx

import (
	"context"
	"time"
)

func ChannelIn[T any](ctx context.Context, ch chan<- T, v T) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case ch <- v:
	}
	return nil
}

func ChannelOut[T any](ctx context.Context, ch <-chan T) (v T, err error) {
	select {
	case <-ctx.Done():
		err = ctx.Err()
		return
	case v = <-ch:
	}
	return
}

func ChannelInWithTimeout[T any](ch chan<- T, v T, timeout time.Duration) error {
	ctx, cf := context.WithTimeout(context.Background(), timeout)
	defer cf()
	return ChannelIn(ctx, ch, v)
}

func ChannelOutWithTimeout[T any](ch <-chan T, timeout time.Duration) (T, error) {
	ctx, cf := context.WithTimeout(context.Background(), timeout)
	defer cf()
	return ChannelOut(ctx, ch)
}
