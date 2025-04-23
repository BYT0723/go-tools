package httpx

import (
	"context"
	"time"
)

var defaultClient = NewClient()

func Get(rawUrl string, ps ...Param) (*Response, error) {
	return defaultClient.Get(context.Background(), rawUrl, ps...)
}

func Getx(ctx context.Context, rawUrl string, ps ...Param) (*Response, error) {
	return defaultClient.Get(ctx, rawUrl, ps...)
}

func GetAny[T any](rawUrl string, ps ...Param) (*Response, *T, error) {
	obj := new(T)
	resp, err := defaultClient.GetAny(context.Background(), rawUrl, obj, ps...)
	return resp, obj, err
}

func GetxAny[T any](ctx context.Context, rawUrl string, ps ...Param) (*Response, *T, error) {
	obj := new(T)
	resp, err := defaultClient.GetAny(ctx, rawUrl, obj, ps...)
	return resp, obj, err
}

func Post(rawUrl string, ps ...Param) (*Response, error) {
	return defaultClient.Post(context.Background(), rawUrl, ps...)
}

func Postx(ctx context.Context, rawUrl string, ps ...Param) (*Response, error) {
	return defaultClient.Post(ctx, rawUrl, ps...)
}

func PostAny[T any](rawUrl string, timeout time.Duration, ps ...Param) (*Response, *T, error) {
	obj := new(T)
	resp, err := defaultClient.PostAny(context.Background(), rawUrl, obj, ps...)
	return resp, obj, err
}

func PostxAny[T any](ctx context.Context, rawUrl string, timeout time.Duration, ps ...Param) (*Response, *T, error) {
	obj := new(T)
	resp, err := defaultClient.PostAny(ctx, rawUrl, obj, ps...)
	return resp, obj, err
}
