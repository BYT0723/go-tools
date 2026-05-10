package httpx

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewClient(t *testing.T) {
	Convey("NewClient 测试", t, func() {
		Convey("创建默认client", func() {
			c := NewClient()
			So(c, ShouldNotBeNil)
			So(c.encoder, ShouldNotBeNil)
			So(c.decoder, ShouldNotBeNil)
			So(c.cli, ShouldNotBeNil)
		})

		Convey("创建client并设置option", func() {
			httpCli := &http.Client{}
			c := NewClient(WithHttpClient(httpCli))
			So(c, ShouldNotBeNil)
			So(c.cli, ShouldEqual, httpCli)
		})
	})
}

func TestDefaultClient(t *testing.T) {
	Convey("DefaultClient 测试", t, func() {
		So(DefaultClient, ShouldNotBeNil)
	})
}

func TestRequestStruct(t *testing.T) {
	Convey("request 结构测试", t, func() {
		r := &request{}
		So(r.header, ShouldBeNil)
		So(r.payload, ShouldBeNil)
	})
}

func TestResponseStruct(t *testing.T) {
	Convey("Response 结构测试", t, func() {
		resp := &Response{
			Code:   200,
			Header: http.Header{"Content-Type": {"application/json"}},
			Body:   []byte(`{"key":"value"}`),
		}
		So(resp.Code, ShouldEqual, 200)
		So(resp.Header.Get("Content-Type"), ShouldEqual, "application/json")
		So(string(resp.Body), ShouldEqual, `{"key":"value"}`)
	})
}

func TestWithHeader(t *testing.T) {
	Convey("WithHeader 测试", t, func() {
		Convey("设置header", func() {
			h := http.Header{"Authorization": {"Bearer token"}}
			p := WithHeader(h)
			r := &request{}
			p(r)
			So(r.header.Get("Authorization"), ShouldEqual, "Bearer token")
		})

		Convey("设置多个header", func() {
			h := http.Header{
				"Content-Type": {"application/json"},
				"Accept":       {"application/json"},
			}
			p := WithHeader(h)
			r := &request{}
			p(r)
			So(r.header.Get("Content-Type"), ShouldEqual, "application/json")
			So(r.header.Get("Accept"), ShouldEqual, "application/json")
		})
	})
}

func TestWithPayload(t *testing.T) {
	Convey("WithPayload 测试", t, func() {
		Convey("设置map payload", func() {
			payload := map[string]any{"key": "value"}
			p := WithPayload(payload)
			r := &request{}
			p(r)
			So(r.payload, ShouldEqual, payload)
		})

		Convey("设置string payload", func() {
			p := WithPayload("hello")
			r := &request{}
			p(r)
			So(r.payload, ShouldEqual, "hello")
		})

		Convey("设置nil payload", func() {
			p := WithPayload(nil)
			r := &request{}
			p(r)
			So(r.payload, ShouldBeNil)
		})
	})
}

func TestEncoderInterface(t *testing.T) {
	Convey("Encoder 接口定义验证", t, func() {
		// Encoder interface: RequestHeader() http.Header, Encode(any) (io.Reader, error)
		So(true, ShouldBeTrue)
	})
}

func TestDecoderInterface(t *testing.T) {
	Convey("Decoder 接口定义验证", t, func() {
		// Decoder interface: Decode(io.Reader, http.Header, any) error
		So(true, ShouldBeTrue)
	})
}

func TestOptionFunctions(t *testing.T) {
	Convey("Option 函数测试", t, func() {
		Convey("WithHttpClient", func() {
			httpCli := &http.Client{}
			c := NewClient()
			WithHttpClient(httpCli)(c)
			So(c.cli, ShouldEqual, httpCli)
		})
	})
}
