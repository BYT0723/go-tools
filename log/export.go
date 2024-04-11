package log

import (
	"github.com/BYT0723/go-tools/log/logger"
	"github.com/BYT0723/go-tools/log/zerologger"
	"go.uber.org/zap"
)

var defaultLogger logger.Logger

func Init(opts ...logger.Option) error {
	zl, err := zerologger.NewInstance(opts...)
	if err != nil {
		return err
	}
	defaultLogger = zl
	return nil
}

func With(kvs ...*logger.Field) logger.Logger {
	return defaultLogger.With(kvs...)
}

func Debug(args ...any) {
	defaultLogger.Debug(args...)
}

func Debugf(format string, args ...any) {
	defaultLogger.Debugf(format, args...)
}

func Info(args ...any) {
	defaultLogger.Info(args...)
}

func Infof(format string, args ...any) {
	defaultLogger.Infof(format, args...)
}

func Warn(args ...any) {
	defaultLogger.Warn(args...)
}

func Warnf(format string, args ...any) {
	defaultLogger.Warnf(format, args...)
}

func Error(args ...any) {
	defaultLogger.Error(args...)
}

func Errorf(format string, args ...any) {
	defaultLogger.Errorf(format, args...)
}

func Panic(args ...any) {
	defaultLogger.Panic(args...)
}

func Panicf(format string, args ...any) {
	defaultLogger.Panicf(format, args...)
}

func Fatal(args ...any) {
	defaultLogger.Fatal(args...)
}

func Fatalf(format string, args ...any) {
	defaultLogger.Fatalf(format, args...)
}

func ZapLogger() (*zap.Logger, bool) {
	return defaultLogger.ZapLogger()
}

func Sync() error {
	return defaultLogger.Sync()
}
