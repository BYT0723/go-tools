package logger

import (
	"go.uber.org/zap"
)

type Logger interface {
	With(kvs ...*Field) Logger
	Debug(args ...any)
	Debugf(format string, args ...any)
	Info(args ...any)
	Infof(format string, args ...any)
	Warn(args ...any)
	Warnf(format string, args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
	Panic(args ...any)
	Panicf(format string, args ...any)
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	ZapLogger() (*zap.Logger, bool)
	Sync() error
}
