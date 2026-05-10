package ginx

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BYT0723/go-tools/logx"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestWithTraceLogger(t *testing.T) {
	t.Run("WithTraceLogger 测试", func(t *testing.T) {
		t.Run("logger为nil不影响请求处理", func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(WithTraceLogger(nil))
			router.GET("/", func(c *gin.Context) {
				c.String(200, "ok")
			})
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			assert.Equal(t, 200, rec.Code)
		})

		t.Run("返回 HandlerFunc 类型", func(t *testing.T) {
			h := WithTraceLogger(nil)
			assert.NotNil(t, h)
		})
	})
}

func TestWithApiLog(t *testing.T) {
	t.Run("WithApiLog 测试", func(t *testing.T) {
		t.Run("生成 HandlerFunc", func(t *testing.T) {
			h := WithApiLog("info")
			assert.NotNil(t, h)
		})

		t.Run("带额外fields的 HandlerFunc", func(t *testing.T) {
			h := WithApiLog("info", func(ctx *gin.Context) []logx.Field {
				return []logx.Field{logx.String("custom", "value")}
			})
			assert.NotNil(t, h)
		})
	})
}

func TestWithTraceID(t *testing.T) {
	t.Run("WithTraceID 测试", func(t *testing.T) {
		t.Run("生成 HandlerFunc", func(t *testing.T) {
			h := WithTraceID(func(ctx *gin.Context) string {
				return "trace-123"
			})
			assert.NotNil(t, h)
		})
	})
}

func TestWithSpanID(t *testing.T) {
	t.Run("WithSpanID 测试", func(t *testing.T) {
		t.Run("生成 HandlerFunc", func(t *testing.T) {
			h := WithSpanID(func(ctx *gin.Context) string { return "span-456" })
			assert.NotNil(t, h)
		})
	})
}

func TestWithRequestID(t *testing.T) {
	t.Run("WithRequestID 测试", func(t *testing.T) {
		t.Run("生成 HandlerFunc", func(t *testing.T) {
			h := WithRequestID(func(ctx *gin.Context) string { return "req-789" })
			assert.NotNil(t, h)
		})
	})
}

func TestWithService(t *testing.T) {
	t.Run("WithService 测试", func(t *testing.T) {
		t.Run("生成 HandlerFunc", func(t *testing.T) {
			h := WithService("my-service")
			assert.NotNil(t, h)
		})
	})
}

func TestWithVersion(t *testing.T) {
	t.Run("WithVersion 测试", func(t *testing.T) {
		t.Run("生成 HandlerFunc", func(t *testing.T) {
			h := WithVersion("1.0.0")
			assert.NotNil(t, h)
		})
	})
}

func TestWithEnv(t *testing.T) {
	t.Run("WithEnv 测试", func(t *testing.T) {
		t.Run("生成 HandlerFunc", func(t *testing.T) {
			h := WithEnv("production")
			assert.NotNil(t, h)
		})
	})
}

func TestWithValue(t *testing.T) {
	t.Run("WithValue 测试", func(t *testing.T) {
		t.Run("生成 HandlerFunc", func(t *testing.T) {
			type keyType struct{}
			h := WithValue(keyType{}, func(ctx *gin.Context) any {
				return "value"
			})
			assert.NotNil(t, h)
		})
	})
}
