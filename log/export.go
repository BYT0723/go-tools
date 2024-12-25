package log

import (
	"fmt"

	"github.com/BYT0723/go-tools/log/logcore"
	"github.com/BYT0723/go-tools/log/zaplogger"
	"github.com/BYT0723/go-tools/log/zerologger"
)

type (
	Logger = logcore.Logger
	Config = logcore.LoggerConf
)

var defaultLogger Logger

func Init(opts ...Option) error {
	logger, err := NewLogger(opts...)
	if err != nil {
		return err
	}
	defaultLogger = logger
	return nil
}

func NewLogger(opts ...Option) (logcore.Logger, error) {
	cfg := &InitConf{LogCfg: logcore.DefaultLoggerConf()}

	for _, opt := range opts {
		opt(cfg)
	}
	switch cfg.Type {
	case TypeZeroLog:
		return zerologger.NewInstance(cfg.LogCfg)
	case TypeZap:
		return zaplogger.NewInstance(cfg.LogCfg)
	default:
		return nil, fmt.Errorf("unknown logger type: %v", cfg.Type)
	}
}

func With(kvs ...Field) logcore.Logger {
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

func Default() Logger {
	return defaultLogger
}

func SetDefault(logger Logger) {
	if logger != nil {
		defaultLogger = logger
	}
}

func Sync() error {
	return defaultLogger.Sync()
}
