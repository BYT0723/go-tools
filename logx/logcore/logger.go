package logcore

type Logger interface {
	With(kvs ...Field) Logger
	Debug(msg string, kvs ...Field)
	Debugf(format string, args ...any)
	Info(msg string, kvs ...Field)
	Infof(format string, args ...any)
	Warn(msg string, kvs ...Field)
	Warnf(format string, args ...any)
	Error(msg string, kvs ...Field)
	Errorf(format string, args ...any)
	Panic(msg string, kvs ...Field)
	Panicf(format string, args ...any)
	Fatal(msg string, kvs ...Field)
	Fatalf(format string, args ...any)
	Log(level string, msg string, kvs ...Field)
	Logf(level string, format string, args ...any)
	Sync() error
	AddCallerSkip(caller int) Logger
}
