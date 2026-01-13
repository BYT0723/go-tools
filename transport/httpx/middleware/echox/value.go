package echox

import (
	"context"

	ctxx "github.com/BYT0723/go-tools/contextx"
	"github.com/labstack/echo/v4"
)

var (
	_ echo.MiddlewareFunc = WithTraceID(func(c echo.Context) string { return "" })
	_ echo.MiddlewareFunc = WithSpanID(func(c echo.Context) string { return "" })
	_ echo.MiddlewareFunc = WithRequestID(func(c echo.Context) string { return "" })
	_ echo.MiddlewareFunc = WithService("")
	_ echo.MiddlewareFunc = WithVersion("")
	_ echo.MiddlewareFunc = WithEnv("")
)

func WithValue(key any, generate func(c echo.Context) any) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.SetRequest(
				c.Request().WithContext(context.WithValue(c.Request().Context(), key, generate(c))),
			)
			return next(c)
		}
	}
}

func WithTraceID(idGenerate func(c echo.Context) string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.SetRequest(
				c.Request().WithContext(ctxx.WithTraceID(c.Request().Context(), idGenerate(c))),
			)
			return next(c)
		}
	}
}

func WithSpanID(idGenerate func(c echo.Context) string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.SetRequest(
				c.Request().WithContext(ctxx.WithSpanID(c.Request().Context(), idGenerate(c))),
			)
			return next(c)
		}
	}
}

func WithRequestID(idGenerate func(c echo.Context) string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.SetRequest(
				c.Request().WithContext(ctxx.WithRequestID(c.Request().Context(), idGenerate(c))),
			)
			return next(c)
		}
	}
}

func WithService(service string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.SetRequest(
				c.Request().WithContext(ctxx.WithService(c.Request().Context(), service)),
			)
			return next(c)
		}
	}
}

func WithVersion(version string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.SetRequest(
				c.Request().WithContext(ctxx.WithVersion(c.Request().Context(), version)),
			)
			return next(c)
		}
	}
}

func WithEnv(env string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.SetRequest(
				c.Request().WithContext(ctxx.WithEnv(c.Request().Context(), env)),
			)
			return next(c)
		}
	}
}
