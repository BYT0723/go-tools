package ginx

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BYT0723/go-tools/logx"
	"github.com/gin-gonic/gin"

	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestWithTraceLogger(t *testing.T) {
	Convey("WithTraceLogger 测试", t, func() {
		Convey("logger为nil不影响请求处理", func() {
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(WithTraceLogger(nil))
			router.GET("/", func(c *gin.Context) {
				c.String(200, "ok")
			})
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			So(rec.Code, ShouldEqual, 200)
		})

		Convey("返回 HandlerFunc 类型", func() {
			h := WithTraceLogger(nil)
			So(h, ShouldNotBeNil)
		})
	})
}

func TestWithApiLog(t *testing.T) {
	Convey("WithApiLog 测试", t, func() {
		Convey("生成 HandlerFunc", func() {
			h := WithApiLog("info")
			So(h, ShouldNotBeNil)
		})

		Convey("带额外fields的 HandlerFunc", func() {
			h := WithApiLog("info", func(ctx *gin.Context) []logx.Field {
				return []logx.Field{logx.String("custom", "value")}
			})
			So(h, ShouldNotBeNil)
		})
	})
}

func TestWithTraceID(t *testing.T) {
	Convey("WithTraceID 测试", t, func() {
		Convey("生成 HandlerFunc", func() {
			h := WithTraceID(func(ctx *gin.Context) string {
				return "trace-123"
			})
			So(h, ShouldNotBeNil)
		})
	})
}

func TestWithSpanID(t *testing.T) {
	Convey("WithSpanID 测试", t, func() {
		Convey("生成 HandlerFunc", func() {
			h := WithSpanID(func(ctx *gin.Context) string { return "span-456" })
			So(h, ShouldNotBeNil)
		})
	})
}

func TestWithRequestID(t *testing.T) {
	Convey("WithRequestID 测试", t, func() {
		Convey("生成 HandlerFunc", func() {
			h := WithRequestID(func(ctx *gin.Context) string { return "req-789" })
			So(h, ShouldNotBeNil)
		})
	})
}

func TestWithService(t *testing.T) {
	Convey("WithService 测试", t, func() {
		Convey("生成 HandlerFunc", func() {
			h := WithService("my-service")
			So(h, ShouldNotBeNil)
		})
	})
}

func TestWithVersion(t *testing.T) {
	Convey("WithVersion 测试", t, func() {
		Convey("生成 HandlerFunc", func() {
			h := WithVersion("1.0.0")
			So(h, ShouldNotBeNil)
		})
	})
}

func TestWithEnv(t *testing.T) {
	Convey("WithEnv 测试", t, func() {
		Convey("生成 HandlerFunc", func() {
			h := WithEnv("production")
			So(h, ShouldNotBeNil)
		})
	})
}

func TestWithValue(t *testing.T) {
	Convey("WithValue 测试", t, func() {
		Convey("生成 HandlerFunc", func() {
			type keyType struct{}
			h := WithValue(keyType{}, func(ctx *gin.Context) any {
				return "value"
			})
			So(h, ShouldNotBeNil)
		})
	})
}
