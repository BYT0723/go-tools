package echox

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BYT0723/go-tools/logx"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestWithTraceLogger(t *testing.T) {
	t.Run("WithTraceLogger 测试", func(t *testing.T) {
		t.Run("logger为nil不影响请求处理", func(t *testing.T) {
			mw := WithTraceLogger(nil)
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			h := mw(func(c echo.Context) error {
				return c.String(200, "ok")
			})
			err := h(c)
			assert.Nil(t, err)
			assert.Equal(t, 200, rec.Code)
		})

		t.Run("返回 MiddlewareFunc 类型", func(t *testing.T) {
			mw := WithTraceLogger(nil)
			assert.NotNil(t, mw)
		})
	})
}

func TestWithApiLog(t *testing.T) {
	t.Run("WithApiLog 测试", func(t *testing.T) {
		t.Run("生成 MiddlewareFunc", func(t *testing.T) {
			mw := WithApiLog("info")
			assert.NotNil(t, mw)
		})

		t.Run("带额外fields的 MiddlewareFunc", func(t *testing.T) {
			mw := WithApiLog("info", func(c echo.Context) []logx.Field {
				return []logx.Field{logx.String("custom", "value")}
			})
			assert.NotNil(t, mw)
		})
	})
}

func TestWithTraceID(t *testing.T) {
	t.Run("WithTraceID 测试", func(t *testing.T) {
		t.Run("设置 TraceID 到 context", func(t *testing.T) {
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
			assert.Nil(t, err)
		})
	})
}

func TestWithSpanID(t *testing.T) {
	t.Run("WithSpanID 测试", func(t *testing.T) {
		t.Run("生成 MiddlewareFunc", func(t *testing.T) {
			mw := WithSpanID(func(c echo.Context) string { return "span-456" })
			assert.NotNil(t, mw)
		})
	})
}

func TestWithRequestID(t *testing.T) {
	t.Run("WithRequestID 测试", func(t *testing.T) {
		t.Run("生成 MiddlewareFunc", func(t *testing.T) {
			mw := WithRequestID(func(c echo.Context) string { return "req-789" })
			assert.NotNil(t, mw)
		})
	})
}

func TestWithService(t *testing.T) {
	t.Run("WithService 测试", func(t *testing.T) {
		t.Run("生成 MiddlewareFunc", func(t *testing.T) {
			mw := WithService("my-service")
			assert.NotNil(t, mw)
		})
	})
}

func TestWithVersion(t *testing.T) {
	t.Run("WithVersion 测试", func(t *testing.T) {
		t.Run("生成 MiddlewareFunc", func(t *testing.T) {
			mw := WithVersion("1.0.0")
			assert.NotNil(t, mw)
		})
	})
}

func TestWithEnv(t *testing.T) {
	t.Run("WithEnv 测试", func(t *testing.T) {
		t.Run("生成 MiddlewareFunc", func(t *testing.T) {
			mw := WithEnv("production")
			assert.NotNil(t, mw)
		})
	})
}

func TestWithValue(t *testing.T) {
	t.Run("WithValue 测试", func(t *testing.T) {
		t.Run("生成 MiddlewareFunc", func(t *testing.T) {
			type keyType struct{}
			mw := WithValue(keyType{}, func(c echo.Context) any {
				return "value"
			})
			assert.NotNil(t, mw)
		})
	})
}
