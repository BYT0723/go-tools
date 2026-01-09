package logx

import (
	"fmt"

	"github.com/BYT0723/go-tools/logx/logcore"
	"github.com/BYT0723/go-tools/logx/zaplogger"
	"github.com/BYT0723/go-tools/logx/zerologger"
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
	return defaultLogger.AddCallerSkip(-1).With(kvs...)
}

func Debug(msg string, kvs ...Field) {
	defaultLogger.Debug(msg, kvs...)
}

func Debugf(format string, args ...any) {
	defaultLogger.Debugf(format, args...)
}

func Info(msg string, kvs ...Field) {
	defaultLogger.Info(msg, kvs...)
}

func Infof(format string, args ...any) {
	defaultLogger.Infof(format, args...)
}

func Warn(msg string, kvs ...Field) {
	defaultLogger.Warn(msg, kvs...)
}

func Warnf(format string, args ...any) {
	defaultLogger.Warnf(format, args...)
}

func Error(msg string, kvs ...Field) {
	defaultLogger.Error(msg, kvs...)
}

func Errorf(format string, args ...any) {
	defaultLogger.Errorf(format, args...)
}

func Panic(msg string, kvs ...Field) {
	defaultLogger.Panic(msg, kvs...)
}

func Panicf(format string, args ...any) {
	defaultLogger.Panicf(format, args...)
}

func Fatal(msg string, kvs ...Field) {
	defaultLogger.Fatal(msg, kvs...)
}

func Fatalf(format string, args ...any) {
	defaultLogger.Fatalf(format, args...)
}

func Default() Logger {
	return defaultLogger.AddCallerSkip(-1)
}

func SetDefault(logger Logger) {
	if logger != nil {
		defaultLogger = logger
	}
}

func Sync() error {
	return defaultLogger.Sync()
}
