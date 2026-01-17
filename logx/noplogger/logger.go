package noplogger

import "github.com/BYT0723/go-tools/logx/logcore"

type NopLogger struct{}

func (l NopLogger) With(kvs ...logcore.Field) logcore.Logger {
	return l
}

func (NopLogger) Debug(msg string, kvs ...logcore.Field) {
}

func (NopLogger) Debugf(format string, args ...any) {
}

func (NopLogger) Info(msg string, kvs ...logcore.Field) {
}

func (NopLogger) Infof(format string, args ...any) {
}

func (NopLogger) Warn(msg string, kvs ...logcore.Field) {
}

func (NopLogger) Warnf(format string, args ...any) {
}

func (NopLogger) Error(msg string, kvs ...logcore.Field) {
}

func (NopLogger) Errorf(format string, args ...any) {
}

func (NopLogger) Panic(msg string, kvs ...logcore.Field) {
}

func (NopLogger) Panicf(format string, args ...any) {
}

func (NopLogger) Fatal(msg string, kvs ...logcore.Field) {
}

func (NopLogger) Fatalf(format string, args ...any) {
}

func (NopLogger) Log(level string, msg string, kvs ...logcore.Field) {
}

func (NopLogger) Logf(level string, format string, args ...any) {
}

func (NopLogger) Sync() error {
	return nil
}

func (l NopLogger) AddCallerSkip(caller int) logcore.Logger {
	return l
}
