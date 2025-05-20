package httpx

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"

	"github.com/BYT0723/go-tools/transport/httpx/decoder"
	"github.com/BYT0723/go-tools/transport/httpx/encoder"
	"golang.org/x/exp/maps"
)

type (
	Client struct {
		encoder Encoder
		decoder Decoder
		cli     *http.Client
	}
	request struct {
		header  http.Header
		payload any
	}
	Response struct {
		Code   int
		Header http.Header
		Body   []byte
	}
	Encoder interface {
		RequestHeader() http.Header
		Encode(any) (io.Reader, error)
	}
	Decoder interface {
		Decode(io.Reader, http.Header, any) error
	}
)

func NewClient(opts ...Option) *Client {
	c := &Client{
		encoder: encoder.JsonEncoder(),
		decoder: decoder.DefaultDecoder(),
		cli:     &http.Client{},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) Get(ctx context.Context, addr string, ps ...Param) (resp *Response, err error) {
	return c.handle(ctx, http.MethodGet, addr, nil, ps...)
}

func (c *Client) Post(ctx context.Context, addr string, ps ...Param) (resp *Response, err error) {
	return c.handle(ctx, http.MethodPost, addr, nil, ps...)
}

func (c *Client) GetAny(ctx context.Context, addr string, result any, ps ...Param) (resp *Response, err error) {
	return c.handle(ctx, http.MethodGet, addr, result, ps...)
}

func (c *Client) PostAny(ctx context.Context, addr string, result any, ps ...Param) (resp *Response, err error) {
	return c.handle(ctx, http.MethodPost, addr, result, ps...)
}

func (c *Client) Do(ctx context.Context, method, addr string, ps ...Param) (resp *Response, err error) {
	return c.handle(ctx, method, addr, nil, ps...)
}

func (c *Client) do(ctx context.Context, method, addr string, ps ...Param) (resp *http.Response, err error) {
	var (
		req   *http.Request
		param = &request{}
	)
	for _, p := range ps {
		p(param)
	}

	if param.payload != nil {
		switch method {
		case http.MethodGet, http.MethodHead, http.MethodDelete:
			var u *url.URL
			u, err = url.Parse(addr)
			if err != nil {
				return
			}
			var (
				query = u.Query()
				t     = reflect.TypeOf(param.payload)
				v     = reflect.ValueOf(param.payload)
			)
			switch t.Kind() {
			case reflect.Struct:
				for i := range t.NumField() {
					value := v.Field(i)
					if value.Kind() == reflect.Interface {
						value = value.Elem()
					}
					fieldName := t.Field(i).Name

					switch value.Kind() {
					case reflect.Slice, reflect.Array:
						for j := range value.Len() {
							query.Add(fieldName, fmt.Sprint(value.Index(j).Interface()))
						}
					default:
						query.Add(fieldName, fmt.Sprint(value.Interface()))
					}
				}
			case reflect.Map:
				if t.Key().Kind() != reflect.String {
					err = errors.New("GET Params Map key type must be string")
					return
				}
				for _, k := range v.MapKeys() {
					value := v.MapIndex(k)
					if value.Kind() == reflect.Interface {
						value = value.Elem()
					}

					switch value.Kind() {
					case reflect.Slice, reflect.Array:
						for j := range value.Len() {
							query.Add(
								fmt.Sprint(k.Interface()),
								fmt.Sprint(value.Index(j).Interface()),
							)
						}
					default:
						query.Add(fmt.Sprint(k.Interface()), fmt.Sprint(value))
					}

				}
			default:
				err = errors.New("GET Params must be Map or Struct")
				return
			}
			u.RawQuery = query.Encode()
			req, err = http.NewRequestWithContext(ctx, method, u.String(), nil)
			if err != nil {
				return
			}
		case http.MethodPost, http.MethodPut, http.MethodPatch:
			var (
				body io.Reader
				eh   bool
			)
			switch v := param.payload.(type) {
			case string:
				body = bytes.NewBufferString(v)
			case []byte:
				body = bytes.NewBuffer(v)
			case io.Reader:
				body = v
			default:
				body, err = c.encoder.Encode(param.payload)
				if err != nil {
					return
				}
				eh = true
			}
			req, err = http.NewRequestWithContext(ctx, method, addr, body)
			if err != nil {
				return
			}
			if eh {
				maps.Copy(req.Header, c.encoder.RequestHeader())
			}
		}
	} else {
		req, err = http.NewRequestWithContext(ctx, method, addr, nil)
		if err != nil {
			return
		}
	}
	if param.header != nil {
		maps.Copy(req.Header, param.header)
	}
	return c.cli.Do(req)
}

// 传入result, 则使用decoder进行解码，respBody将会返回空值
// 反之result为nil, respBody将会返回原始数据
func (c *Client) handle(ctx context.Context, method, addr string, result any, ps ...Param) (resp *Response, err error) {
	rp, err := c.do(ctx, method, addr, ps...)
	if err != nil {
		return
	}
	defer rp.Body.Close()

	resp = new(Response)

	resp.Code = rp.StatusCode
	resp.Header = rp.Header

	resp.Body, err = io.ReadAll(rp.Body)
	if err != nil {
		return
	}

	if result != nil {
		if c.decoder == nil {
			err = errors.New("decoder is nil")
			return
		}
		err = c.decoder.Decode(bytes.NewBuffer(resp.Body), rp.Header, result)
		if err != nil {
			err = fmt.Errorf("response decode err: %v, source: \"%s\"", err, resp.Body)
		}
	}
	return
}
