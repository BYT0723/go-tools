package ctxx

import (
	"context"
	"net"
	"sync"
	"testing"

	"github.com/BYT0723/go-tools/logx/noplogger"

	"github.com/stretchr/testify/assert"
)

func TestErrorVars(t *testing.T) {
	t.Run("错误变量测试", func(t *testing.T) {
		assert.Equal(t, "api key not found", ErrApiKeyNotFound.Error())
		assert.Equal(t, "waitgroup not found", ErrWaitGroupNotFound.Error())
		assert.Equal(t, "listener not found", ErrListenerNotFound.Error())
		assert.Equal(t, "logger not found", ErrLoggerNotFound.Error())
	})
}

func TestListener(t *testing.T) {
	t.Run("Listener 测试", func(t *testing.T) {
		t.Run("获取不存在的Listener返回错误", func(t *testing.T) {
			ctx := context.Background()
			_, err := Listener(ctx)
			assert.Equal(t, ErrListenerNotFound, err)
		})

		t.Run("设置后获取Listener成功", func(t *testing.T) {
			l, _ := net.Listen("tcp", "127.0.0.1:0")
			defer l.Close()

			ctx := WithListener(context.Background(), l)
			result, err := Listener(ctx)
			assert.Nil(t, err)
			assert.Equal(t, l, result)
		})
	})
}

func TestLogger(t *testing.T) {
	t.Run("Logger 测试", func(t *testing.T) {
		t.Run("获取不存在的Logger返回错误", func(t *testing.T) {
			ctx := context.Background()
			_, err := Logger(ctx)
			assert.Equal(t, ErrLoggerNotFound, err)
		})

		t.Run("设置后获取Logger成功", func(t *testing.T) {
			logger := &noplogger.NopLogger{}
			ctx := WithLogger(context.Background(), logger)
			result, err := Logger(ctx)
			assert.Nil(t, err)
			assert.Equal(t, logger, result)
		})
	})
}

func TestTraceID(t *testing.T) {
	t.Run("TraceID 测试", func(t *testing.T) {
		t.Run("获取不存在的TraceID返回错误", func(t *testing.T) {
			_, err := TraceID(context.Background())
			assert.Equal(t, ErrApiKeyNotFound, err)
		})

		t.Run("设置后获取TraceID成功", func(t *testing.T) {
			ctx := WithTraceID(context.Background(), "trace-123")
			id, err := TraceID(ctx)
			assert.Nil(t, err)
			assert.Equal(t, "trace-123", id)
		})
	})
}

func TestSpanID(t *testing.T) {
	t.Run("SpanID 测试", func(t *testing.T) {
		t.Run("获取不存在的SpanID返回错误", func(t *testing.T) {
			_, err := SpanID(context.Background())
			assert.Equal(t, ErrApiKeyNotFound, err)
		})

		t.Run("设置后获取SpanID成功", func(t *testing.T) {
			ctx := WithSpanID(context.Background(), "span-456")
			id, err := SpanID(ctx)
			assert.Nil(t, err)
			assert.Equal(t, "span-456", id)
		})
	})
}

func TestRequestID(t *testing.T) {
	t.Run("RequestID 测试", func(t *testing.T) {
		t.Run("获取不存在的RequestID返回错误", func(t *testing.T) {
			_, err := RequestID(context.Background())
			assert.Equal(t, ErrApiKeyNotFound, err)
		})

		t.Run("设置后获取RequestID成功", func(t *testing.T) {
			ctx := WithRequestID(context.Background(), "req-789")
			id, err := RequestID(ctx)
			assert.Nil(t, err)
			assert.Equal(t, "req-789", id)
		})
	})
}

func TestService(t *testing.T) {
	t.Run("Service 测试", func(t *testing.T) {
		t.Run("获取不存在的Service返回错误", func(t *testing.T) {
			_, err := Service(context.Background())
			assert.Equal(t, ErrApiKeyNotFound, err)
		})

		t.Run("设置后获取Service成功", func(t *testing.T) {
			ctx := WithService(context.Background(), "my-service")
			s, err := Service(ctx)
			assert.Nil(t, err)
			assert.Equal(t, "my-service", s)
		})
	})
}

func TestVersion(t *testing.T) {
	t.Run("Version 测试", func(t *testing.T) {
		t.Run("获取不存在的Version返回错误", func(t *testing.T) {
			_, err := Version(context.Background())
			assert.Equal(t, ErrApiKeyNotFound, err)
		})

		t.Run("设置后获取Version成功", func(t *testing.T) {
			ctx := WithVersion(context.Background(), "1.0.0")
			v, err := Version(ctx)
			assert.Nil(t, err)
			assert.Equal(t, "1.0.0", v)
		})
	})
}

func TestEnv(t *testing.T) {
	t.Run("Env 测试", func(t *testing.T) {
		t.Run("获取不存在的Env返回错误", func(t *testing.T) {
			_, err := Env(context.Background())
			assert.Equal(t, ErrApiKeyNotFound, err)
		})

		t.Run("设置后获取Env成功", func(t *testing.T) {
			ctx := WithEnv(context.Background(), "production")
			e, err := Env(ctx)
			assert.Nil(t, err)
			assert.Equal(t, "production", e)
		})
	})
}

func TestWithWaitGroup(t *testing.T) {
	t.Run("WaitGroup 测试", func(t *testing.T) {
		t.Run("获取不存在的WaitGroup返回错误", func(t *testing.T) {
			_, err := WaitGroup(context.Background())
			assert.Equal(t, ErrWaitGroupNotFound, err)
		})

		t.Run("设置后获取WaitGroup成功", func(t *testing.T) {
			var wg sync.WaitGroup
			ctx := WithWaitGroup(context.Background(), &wg)
			result, err := WaitGroup(ctx)
			assert.Nil(t, err)
			assert.Equal(t, &wg, result)
		})
	})
}
