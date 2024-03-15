package http

import (
	"context"
	"net/http"
	"time"
)

var DefaultClient = NewClient()

func Get(rawUrl string, header http.Header, data any) (code int, body []byte, err error) {
	return DefaultClient.Get(context.Background(), rawUrl, header, data)
}

func GetWithContext(ctx context.Context, rawUrl string, header http.Header, data any) (code int, body []byte, err error) {
	return DefaultClient.Get(ctx, rawUrl, header, data)
}

func GetTimeout(rawUrl string, header http.Header, data any, timeout time.Duration) (code int, resp []byte, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return DefaultClient.Get(ctx, rawUrl, header, data)
}

func GetAny[T any](rawUrl string, header http.Header, data any) (code int, resp *T, err error) {
	resp = new(T)
	code, err = DefaultClient.GetAny(context.Background(), rawUrl, header, data, resp)
	return
}

func GetAnyWithContext[T any](ctx context.Context, rawUrl string, header http.Header, data any) (code int, resp *T, err error) {
	resp = new(T)
	code, err = DefaultClient.GetAny(ctx, rawUrl, header, data, resp)
	return
}

func GetAnyTimeout[T any](rawUrl string, header http.Header, data any, timeout time.Duration) (code int, resp *T, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp = new(T)
	code, err = DefaultClient.GetAny(ctx, rawUrl, header, data, resp)
	return
}

func Post(rawUrl string, header http.Header, data any) (code int, body []byte, err error) {
	return DefaultClient.Post(context.Background(), rawUrl, header, data)
}

func PostWithContext(ctx context.Context, rawUrl string, header http.Header, data any) (code int, body []byte, err error) {
	return DefaultClient.Post(ctx, rawUrl, header, data)
}

func PostTimeout(rawUrl string, header http.Header, data any, timeout time.Duration) (code int, body []byte, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return DefaultClient.Post(ctx, rawUrl, header, data)
}

func PostAny[T any](rawUrl string, header http.Header, data any) (code int, resp *T, err error) {
	resp = new(T)
	code, err = DefaultClient.PostAny(context.Background(), rawUrl, header, data, resp)
	return
}

func PostAnyWithContext[T any](ctx context.Context, rawUrl string, header http.Header, data any, timeout time.Duration) (code int, resp *T, err error) {
	resp = new(T)
	code, err = DefaultClient.PostAny(ctx, rawUrl, header, data, resp)
	return
}

func PostAnyTimeout[T any](rawUrl string, header http.Header, data any, timeout time.Duration) (code int, resp *T, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp = new(T)
	code, err = DefaultClient.PostAny(ctx, rawUrl, header, data, resp)
	return
}
