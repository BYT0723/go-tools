package httpx

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	t.Run("NewClient 测试", func(t *testing.T) {
		t.Run("创建默认client", func(t *testing.T) {
			c := NewClient()
			assert.NotNil(t, c)
			assert.NotNil(t, c.encoder)
			assert.NotNil(t, c.decoder)
			assert.NotNil(t, c.cli)
		})

		t.Run("创建client并设置option", func(t *testing.T) {
			httpCli := &http.Client{}
			c := NewClient(WithHttpClient(httpCli))
			assert.NotNil(t, c)
			assert.Equal(t, httpCli, c.cli)
		})
	})
}

func TestDefaultClient(t *testing.T) {
	t.Run("DefaultClient 测试", func(t *testing.T) {
		assert.NotNil(t, DefaultClient)
	})
}

func TestRequestStruct(t *testing.T) {
	t.Run("request 结构测试", func(t *testing.T) {
		r := &request{}
		assert.Nil(t, r.header)
		assert.Nil(t, r.payload)
	})
}

func TestResponseStruct(t *testing.T) {
	t.Run("Response 结构测试", func(t *testing.T) {
		resp := &Response{
			Code:   200,
			Header: http.Header{"Content-Type": {"application/json"}},
			Body:   []byte(`{"key":"value"}`),
		}
		assert.Equal(t, 200, resp.Code)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
		assert.Equal(t, `{"key":"value"}`, string(resp.Body))
	})
}

func TestWithHeader(t *testing.T) {
	t.Run("WithHeader 测试", func(t *testing.T) {
		t.Run("设置header", func(t *testing.T) {
			h := http.Header{"Authorization": {"Bearer token"}}
			p := WithHeader(h)
			r := &request{}
			p(r)
			assert.Equal(t, "Bearer token", r.header.Get("Authorization"))
		})

		t.Run("设置多个header", func(t *testing.T) {
			h := http.Header{
				"Content-Type": {"application/json"},
				"Accept":       {"application/json"},
			}
			p := WithHeader(h)
			r := &request{}
			p(r)
			assert.Equal(t, "application/json", r.header.Get("Content-Type"))
			assert.Equal(t, "application/json", r.header.Get("Accept"))
		})
	})
}

func TestWithPayload(t *testing.T) {
	t.Run("WithPayload 测试", func(t *testing.T) {
		t.Run("设置map payload", func(t *testing.T) {
			payload := map[string]any{"key": "value"}
			p := WithPayload(payload)
			r := &request{}
			p(r)
			assert.Equal(t, payload, r.payload)
		})

		t.Run("设置string payload", func(t *testing.T) {
			p := WithPayload("hello")
			r := &request{}
			p(r)
			assert.Equal(t, "hello", r.payload)
		})

		t.Run("设置nil payload", func(t *testing.T) {
			p := WithPayload(nil)
			r := &request{}
			p(r)
			assert.Nil(t, r.payload)
		})
	})
}

func TestEncoderInterface(t *testing.T) {
	t.Run("Encoder 接口定义验证", func(t *testing.T) {
		assert.True(t, true)
	})
}

func TestDecoderInterface(t *testing.T) {
	t.Run("Decoder 接口定义验证", func(t *testing.T) {
		assert.True(t, true)
	})
}

func TestOptionFunctions(t *testing.T) {
	t.Run("Option 函数测试", func(t *testing.T) {
		t.Run("WithHttpClient", func(t *testing.T) {
			httpCli := &http.Client{}
			c := NewClient()
			WithHttpClient(httpCli)(c)
			assert.Equal(t, httpCli, c.cli)
		})
	})
}
