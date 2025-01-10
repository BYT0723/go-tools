package logcore

type Logger interface {
	With(kvs ...Field) Logger
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
	Log(level string, args ...any)
	Logf(level string, format string, args ...any)
	Sync() error
}
