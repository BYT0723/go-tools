package ctxx

import (
	"context"
	"net"
	"sync"
	"testing"

	"github.com/BYT0723/go-tools/logx/noplogger"

	. "github.com/smartystreets/goconvey/convey"
)

func TestErrorVars(t *testing.T) {
	Convey("错误变量测试", t, func() {
		So(ErrApiKeyNotFound.Error(), ShouldEqual, "api key not found")
		So(ErrWaitGroupNotFound.Error(), ShouldEqual, "waitgroup not found")
		So(ErrListenerNotFound.Error(), ShouldEqual, "listener not found")
		So(ErrLoggerNotFound.Error(), ShouldEqual, "logger not found")
	})
}

func TestListener(t *testing.T) {
	Convey("Listener 测试", t, func() {
		Convey("获取不存在的Listener返回错误", func() {
			ctx := context.Background()
			_, err := Listener(ctx)
			So(err, ShouldEqual, ErrListenerNotFound)
		})

		Convey("设置后获取Listener成功", func() {
			l, _ := net.Listen("tcp", "127.0.0.1:0")
			defer l.Close()

			ctx := WithListener(context.Background(), l)
			result, err := Listener(ctx)
			So(err, ShouldBeNil)
			So(result, ShouldEqual, l)
		})
	})
}

func TestLogger(t *testing.T) {
	Convey("Logger 测试", t, func() {
		Convey("获取不存在的Logger返回错误", func() {
			ctx := context.Background()
			_, err := Logger(ctx)
			So(err, ShouldEqual, ErrLoggerNotFound)
		})

		Convey("设置后获取Logger成功", func() {
			logger := &noplogger.NopLogger{}
			ctx := WithLogger(context.Background(), logger)
			result, err := Logger(ctx)
			So(err, ShouldBeNil)
			So(result, ShouldEqual, logger)
		})
	})
}

func TestTraceID(t *testing.T) {
	Convey("TraceID 测试", t, func() {
		Convey("获取不存在的TraceID返回错误", func() {
			_, err := TraceID(context.Background())
			So(err, ShouldEqual, ErrApiKeyNotFound)
		})

		Convey("设置后获取TraceID成功", func() {
			ctx := WithTraceID(context.Background(), "trace-123")
			id, err := TraceID(ctx)
			So(err, ShouldBeNil)
			So(id, ShouldEqual, "trace-123")
		})
	})
}

func TestSpanID(t *testing.T) {
	Convey("SpanID 测试", t, func() {
		Convey("获取不存在的SpanID返回错误", func() {
			_, err := SpanID(context.Background())
			So(err, ShouldEqual, ErrApiKeyNotFound)
		})

		Convey("设置后获取SpanID成功", func() {
			ctx := WithSpanID(context.Background(), "span-456")
			id, err := SpanID(ctx)
			So(err, ShouldBeNil)
			So(id, ShouldEqual, "span-456")
		})
	})
}

func TestRequestID(t *testing.T) {
	Convey("RequestID 测试", t, func() {
		Convey("获取不存在的RequestID返回错误", func() {
			_, err := RequestID(context.Background())
			So(err, ShouldEqual, ErrApiKeyNotFound)
		})

		Convey("设置后获取RequestID成功", func() {
			ctx := WithRequestID(context.Background(), "req-789")
			id, err := RequestID(ctx)
			So(err, ShouldBeNil)
			So(id, ShouldEqual, "req-789")
		})
	})
}

func TestService(t *testing.T) {
	Convey("Service 测试", t, func() {
		Convey("获取不存在的Service返回错误", func() {
			_, err := Service(context.Background())
			So(err, ShouldEqual, ErrApiKeyNotFound)
		})

		Convey("设置后获取Service成功", func() {
			ctx := WithService(context.Background(), "my-service")
			s, err := Service(ctx)
			So(err, ShouldBeNil)
			So(s, ShouldEqual, "my-service")
		})
	})
}

func TestVersion(t *testing.T) {
	Convey("Version 测试", t, func() {
		Convey("获取不存在的Version返回错误", func() {
			_, err := Version(context.Background())
			So(err, ShouldEqual, ErrApiKeyNotFound)
		})

		Convey("设置后获取Version成功", func() {
			ctx := WithVersion(context.Background(), "1.0.0")
			v, err := Version(ctx)
			So(err, ShouldBeNil)
			So(v, ShouldEqual, "1.0.0")
		})
	})
}

func TestEnv(t *testing.T) {
	Convey("Env 测试", t, func() {
		Convey("获取不存在的Env返回错误", func() {
			_, err := Env(context.Background())
			So(err, ShouldEqual, ErrApiKeyNotFound)
		})

		Convey("设置后获取Env成功", func() {
			ctx := WithEnv(context.Background(), "production")
			e, err := Env(ctx)
			So(err, ShouldBeNil)
			So(e, ShouldEqual, "production")
		})
	})
}

func TestWithWaitGroup(t *testing.T) {
	Convey("WaitGroup 测试", t, func() {
		Convey("获取不存在的WaitGroup返回错误", func() {
			_, err := WaitGroup(context.Background())
			So(err, ShouldEqual, ErrWaitGroupNotFound)
		})

		Convey("设置后获取WaitGroup成功", func() {
			var wg sync.WaitGroup
			ctx := WithWaitGroup(context.Background(), &wg)
			result, err := WaitGroup(ctx)
			So(err, ShouldBeNil)
			So(result, ShouldEqual, &wg)
		})
	})
}
