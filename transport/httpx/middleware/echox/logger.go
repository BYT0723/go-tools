package echox

import (
	"time"

	ctxx "github.com/BYT0723/go-tools/contextx"
	"github.com/BYT0723/go-tools/logx"
	"github.com/labstack/echo/v4"
)

var (
	_ echo.MiddlewareFunc = WithTraceLogger(nil)
	_ echo.MiddlewareFunc = WithApiLog("info")
)

func WithApiLog(level string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var (
				start = time.Now()
				err   = next(c)
			)

			l, lerr := ctxx.Logger(c.Request().Context())
			if lerr != nil {
				return err
			}

			var (
				// 字段
				fields = make([]logx.Field, 0, 10)
				// 结束时间
				latency = time.Since(start)
			)

			if traceID, err := ctxx.TraceID(c.Request().Context()); err == nil {
				fields = append(fields, logx.String("trace_id", traceID))
			}
			if spanID, err := ctxx.SpanID(c.Request().Context()); err == nil {
				fields = append(fields, logx.String("span_id", spanID))
			}
			if reqID, err := ctxx.RequestID(c.Request().Context()); err == nil {
				fields = append(fields, logx.String("request_id", reqID))
			}
			if service, err := ctxx.Service(c.Request().Context()); err == nil {
				fields = append(fields, logx.String("service", service))
			}
			if version, err := ctxx.Version(c.Request().Context()); err == nil {
				fields = append(fields, logx.String("version", version))
			}
			if env, err := ctxx.Env(c.Request().Context()); err == nil {
				fields = append(fields, logx.String("env", env))
			}

			fields = append(fields,
				logx.String("method", c.Request().Method),
				logx.String("path", c.Request().URL.Path),
				logx.Int("status", c.Response().Status),
				logx.Duration("latency", latency),
			)
			l.Log(level, "API REQUEST", fields...)

			return err
		}
	}
}

func WithTraceLogger(logger logx.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if logger != nil {
				c.SetRequest(
					c.Request().WithContext(ctxx.WithLogger(c.Request().Context(), logger)),
				)
			}
			return next(c)
		}
	}
}
