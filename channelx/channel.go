package channelx

import (
	"context"
	"time"
)

func TryIn[T any](ch chan<- T, v T) bool {
	select {
	case ch <- v:
		return true
	default:
		return false
	}
}

func TryOut[T any](ch <-chan T) (v T, ok bool) {
	select {
	case v = <-ch:
		return v, true
	default:
		return v, false
	}
}

func In[T any](ctx context.Context, ch chan<- T, v T) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case ch <- v:
	}
	return nil
}

func Out[T any](ctx context.Context, ch <-chan T) (v T, err error) {
	select {
	case <-ctx.Done():
		err = ctx.Err()
		return
	case v = <-ch:
	}
	return
}

func InTimeout[T any](ch chan<- T, v T, timeout time.Duration) error {
	ctx, cf := context.WithTimeout(context.Background(), timeout)
	defer cf()
	return In(ctx, ch, v)
}

func OutTimeout[T any](ch <-chan T, timeout time.Duration) (T, error) {
	ctx, cf := context.WithTimeout(context.Background(), timeout)
	defer cf()
	return Out(ctx, ch)
}

func InDeadline[T any](ch chan<- T, v T, t time.Time) error {
	ctx, cf := context.WithDeadline(context.Background(), t)
	defer cf()
	return In(ctx, ch, v)
}

func OutDeadline[T any](ch <-chan T, t time.Time) (T, error) {
	ctx, cf := context.WithDeadline(context.Background(), t)
	defer cf()
	return Out(ctx, ch)
}
