package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
)

type Client struct {
	encoder     Encoder
	decoder     Decoder
	innerClient *http.Client
}

func NewClient(opts ...Option) *Client {
	c := &Client{
		encoder:     JsonEncoder,
		decoder:     JsonDecoder,
		innerClient: &http.Client{},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) Get(ctx context.Context, rawUrl string, header http.Header, payload any) (code int, body []byte, err error) {
	return c.handle(ctx, http.MethodGet, rawUrl, payload, header, nil, false)
}

func (c *Client) Post(ctx context.Context, rawUrl string, header http.Header, payload any) (code int, body []byte, err error) {
	return c.handle(ctx, http.MethodPost, rawUrl, payload, header, nil, false)
}

func (c *Client) GetAny(ctx context.Context, rawUrl string, header http.Header, payload any, result any) (code int, err error) {
	code, _, err = c.handle(ctx, http.MethodGet, rawUrl, payload, header, result, true)
	return
}

func (c *Client) PostAny(ctx context.Context, rawUrl string, header http.Header, payload any, result any) (code int, err error) {
	code, _, err = c.handle(ctx, http.MethodPost, rawUrl, payload, header, result, true)
	return
}

func (c *Client) Do(ctx context.Context, method, rawUrl string, header http.Header, payload any) (code int, body io.ReadCloser, err error) {
	var buf bytes.Buffer
	defer buf.Reset()

	if payload != nil {
		if method == http.MethodGet {
			u, err := url.Parse(rawUrl)
			if err != nil {
				return 0, nil, err
			}
			var (
				query = u.Query()
				t     = reflect.TypeOf(payload)
				v     = reflect.ValueOf(payload)
			)
			switch t.Kind() {
			case reflect.Struct:
				for i := 0; i < t.NumField(); i++ {
					query.Add(t.Field(i).Name, fmt.Sprint(v.Field(i).Interface()))
				}
			case reflect.Map:
				if t.Key().Kind() != reflect.String {
					return 0, nil, errors.New("GET Params Map key type must be string")
				}
				for _, k := range v.MapKeys() {
					query.Add(fmt.Sprint(k.Interface()), fmt.Sprint(v.MapIndex(k)))
				}
			default:
				return 0, nil, errors.New("GET Params need [Map|Struct]")
			}
			u.RawQuery = query.Encode()
			rawUrl = u.String()
		} else {
			bs, err := c.encoder(ctx, payload)
			if err != nil {
				return 0, nil, err
			}
			buf.Write(bs)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, rawUrl, &buf)
	if err != nil {
		return
	}

	for k, vs := range header {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	resp, err := c.innerClient.Do(req)
	if err != nil {
		return
	}

	code = resp.StatusCode
	body = resp.Body
	return
}

func (c *Client) handle(ctx context.Context, method, rawUrl string, payload any, header http.Header, result any, isDecode bool) (code int, body []byte, err error) {
	code, resp, err := c.Do(ctx, method, rawUrl, header, payload)
	if err != nil {
		return
	}
	defer resp.Close()

	body, err = io.ReadAll(resp)
	if err != nil {
		return
	}

	if isDecode {
		err = c.decoder(ctx, bytes.NewBuffer(body), result)
		if err != nil {
			err = fmt.Errorf("response decode err: %v, source: \"%s\"", err, body)
		}
	}
	return
}
