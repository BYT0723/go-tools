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
	_ echo.MiddlewareFunc = WithApiKey(func(c echo.Context) string { return "" })
)

func WithApiLog(level string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var (
				l, err1   = ctxx.Logger(c.Request().Context())
				apiKey, _ = ctxx.ApiKey(c.Request().Context())
				start     = time.Now()
			)
			err := next(c)
			if err1 != nil {
				return err
			}
			// 结束时间
			latency := time.Since(start)

			l.With(
				logx.String("method", c.Request().Method),
				logx.String("path", c.Request().URL.Path),
				logx.Int("status", c.Response().Status),
				logx.String("api_key", apiKey),
				logx.Duration("latency", latency),
			).Log(level, "API Request")

			return err
		}
	}
}

func WithTraceLogger(logger logx.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.SetRequest(c.Request().WithContext(ctxx.WithLogger(c.Request().Context(), logger)))
			return next(c)
		}
	}
}

func WithApiKey(keyGenerate func(c echo.Context) string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.SetRequest(
				c.Request().WithContext(ctxx.WithApiKey(c.Request().Context(), keyGenerate(c))),
			)
			return next(c)
		}
	}
}
