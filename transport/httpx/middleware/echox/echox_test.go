package echox

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BYT0723/go-tools/logx"
	"github.com/labstack/echo/v4"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWithTraceLogger(t *testing.T) {
	Convey("WithTraceLogger 测试", t, func() {
		Convey("logger为nil不影响请求处理", func() {
			mw := WithTraceLogger(nil)
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			h := mw(func(c echo.Context) error {
				return c.String(200, "ok")
			})
			err := h(c)
			So(err, ShouldBeNil)
			So(rec.Code, ShouldEqual, 200)
		})

		Convey("返回 MiddlewareFunc 类型", func() {
			mw := WithTraceLogger(nil)
			So(mw, ShouldNotBeNil)
		})
	})
}

func TestWithApiLog(t *testing.T) {
	Convey("WithApiLog 测试", t, func() {
		Convey("生成 MiddlewareFunc", func() {
			mw := WithApiLog("info")
			So(mw, ShouldNotBeNil)
		})

		Convey("带额外fields的 MiddlewareFunc", func() {
			mw := WithApiLog("info", func(c echo.Context) []logx.Field {
				return []logx.Field{logx.String("custom", "value")}
			})
			So(mw, ShouldNotBeNil)
		})
	})
}

func TestWithTraceID(t *testing.T) {
	Convey("WithTraceID 测试", t, func() {
		Convey("设置 TraceID 到 context", func() {
			mw := WithTraceID(func(c echo.Context) string {
				return "trace-123"
			})
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			h := mw(func(c echo.Context) error {
				return c.String(200, "ok")
			})
			err := h(c)
			So(err, ShouldBeNil)
		})
	})
}

func TestWithSpanID(t *testing.T) {
	Convey("WithSpanID 测试", t, func() {
		Convey("生成 MiddlewareFunc", func() {
			mw := WithSpanID(func(c echo.Context) string { return "span-456" })
			So(mw, ShouldNotBeNil)
		})
	})
}

func TestWithRequestID(t *testing.T) {
	Convey("WithRequestID 测试", t, func() {
		Convey("生成 MiddlewareFunc", func() {
			mw := WithRequestID(func(c echo.Context) string { return "req-789" })
			So(mw, ShouldNotBeNil)
		})
	})
}

func TestWithService(t *testing.T) {
	Convey("WithService 测试", t, func() {
		Convey("生成 MiddlewareFunc", func() {
			mw := WithService("my-service")
			So(mw, ShouldNotBeNil)
		})
	})
}

func TestWithVersion(t *testing.T) {
	Convey("WithVersion 测试", t, func() {
		Convey("生成 MiddlewareFunc", func() {
			mw := WithVersion("1.0.0")
			So(mw, ShouldNotBeNil)
		})
	})
}

func TestWithEnv(t *testing.T) {
	Convey("WithEnv 测试", t, func() {
		Convey("生成 MiddlewareFunc", func() {
			mw := WithEnv("production")
			So(mw, ShouldNotBeNil)
		})
	})
}

func TestWithValue(t *testing.T) {
	Convey("WithValue 测试", t, func() {
		Convey("生成 MiddlewareFunc", func() {
			type keyType struct{}
			mw := WithValue(keyType{}, func(c echo.Context) any {
				return "value"
			})
			So(mw, ShouldNotBeNil)
		})
	})
}
