package httpx

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
)

var DefaultClient = NewClient()

// Get sends a GET request using DefaultClient to the given rawUrl.
//
// The optional params (ps) may include headers, payload, etc.
func Get(rawUrl string, ps ...Param) (*Response, error) {
	return DefaultClient.Get(context.Background(), rawUrl, ps...)
}

// Getx is the context-aware version of Get.
func Getx(ctx context.Context, rawUrl string, ps ...Param) (*Response, error) {
	return DefaultClient.Get(ctx, rawUrl, ps...)
}

// GetAny sends a GET request using DefaultClient and decodes the response
// body into a value of type T using the Client's decoder.
func GetAny[T any](rawUrl string, ps ...Param) (*Response, T, error) {
	return handleAny[T](context.Background(), http.MethodGet, rawUrl, ps...)
}

// GetxAny is the context-aware version of GetAny.
func GetxAny[T any](ctx context.Context, rawUrl string, ps ...Param) (*Response, T, error) {
	return handleAny[T](ctx, http.MethodGet, rawUrl, ps...)
}

// Post sends a POST request using DefaultClient to the given rawUrl.
func Post(rawUrl string, ps ...Param) (*Response, error) {
	return DefaultClient.Post(context.Background(), rawUrl, ps...)
}

// Postx is the context-aware version of Post.
func Postx(ctx context.Context, rawUrl string, ps ...Param) (*Response, error) {
	return DefaultClient.Post(ctx, rawUrl, ps...)
}

// PostAny sends a POST request using DefaultClient and decodes the response
// body into a value of type T using the Client's decoder.
func PostAny[T any](rawUrl string, ps ...Param) (*Response, T, error) {
	return handleAny[T](context.Background(), http.MethodPost, rawUrl, ps...)
}

// PostxAny is the context-aware version of PostAny.
func PostxAny[T any](ctx context.Context, rawUrl string, ps ...Param) (*Response, T, error) {
	return handleAny[T](ctx, http.MethodPost, rawUrl, ps...)
}

func handleAny[T any](ctx context.Context, method, rawUrl string, ps ...Param) (resp *Response, obj T, err error) {
	t := reflect.TypeOf(obj)

	switch t.Kind() {
	case reflect.Map:
		obj = reflect.MakeMap(t).Interface().(T)
	case reflect.Ptr:
		et := t.Elem()
		if et.Kind() == reflect.Struct {
			obj = reflect.New(et).Interface().(T)
		} else {
			return nil, obj, fmt.Errorf("expected a pointer to struct, but got %s", et.Kind())
		}
		resp, err = DefaultClient.handle(ctx, method, rawUrl, obj, ps...)
		return
	case reflect.Chan, reflect.Func:
		return nil, obj, fmt.Errorf("%s are not allowed", t.Kind())
	}

	resp, err = DefaultClient.handle(ctx, method, rawUrl, &obj, ps...)
	return
}

// Download downloads the file from the given url to the given filepath.
func Download(url string, filepath string) (err error) {
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return
	}
	defer f.Close()

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code error: %d", resp.StatusCode)
		return
	}
	_, err = io.Copy(f, resp.Body)
	return
}
