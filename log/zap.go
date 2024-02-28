package log

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
	Cfg    *LoggerConf
}

type LoggerConf struct {
	Filename string
	Level    string
	// uint: MB
	MaxBackups int
	MaxSize    int
	// uint: DAY
	MaxAge int
}

var (
	l    *Logger
	once sync.Once
)

func Init(opts ...Option) {
	once.Do(func() {
		l = &Logger{
			Cfg: &LoggerConf{
				Filename:   "logs/app.log",
				Level:      "info",
				MaxBackups: 3,
				MaxSize:    10,
				MaxAge:     7,
			},
		}

		for _, opt := range opts {
			opt(l)
		}

		level, err := zap.ParseAtomicLevel(l.Cfg.Level)
		if err != nil {
			panic(err)
		}

		var (
			console = newConsoleCore(level)
			debug   = newCore(zap.NewAtomicLevelAt(zap.DebugLevel), l.Cfg)
			info    = newCore(zap.NewAtomicLevelAt(zap.InfoLevel), l.Cfg)
			warn    = newCore(zap.NewAtomicLevelAt(zap.WarnLevel), l.Cfg)
			error   = newCore(zap.NewAtomicLevelAt(zap.ErrorLevel), l.Cfg)
			dpanic  = newCore(zap.NewAtomicLevelAt(zap.DPanicLevel), l.Cfg)
		)

		core := zapcore.NewTee(console, debug, info, warn, error, dpanic)
		l.logger = zap.New(core)
	})
}

func Debug(args ...any) {
	l.logger.Sugar().Debug(args)
}

func Debugf(msg string, args ...any) {
	l.logger.Sugar().Debugf(msg, args)
}

func Debugw(msg string, keyValues ...any) {
	l.logger.Sugar().Debugw(msg, keyValues)
}

func Info(args ...any) {
	l.logger.Sugar().Info(args)
}

func Infof(msg string, args ...any) {
	l.logger.Sugar().Infof(msg, args)
}

func Infow(msg string, keyValues ...any) {
	l.logger.Sugar().Infow(msg, keyValues)
}

func Warn(args ...any) {
	l.logger.Sugar().Warn(args)
}

func Warnf(msg string, args ...any) {
	l.logger.Sugar().Warnf(msg, args)
}

func Warnw(msg string, keyValues ...any) {
	l.logger.Sugar().Warnw(msg, keyValues)
}

func Error(args ...any) {
	l.logger.Sugar().Error(args)
}

func Errorf(msg string, args ...any) {
	l.logger.Sugar().Errorf(msg, args)
}

func Errorw(msg string, keyValues ...any) {
	l.logger.Sugar().Errorw(msg, keyValues)
}

func Fatal(args ...any) {
	l.logger.Sugar().Fatal(args)
}

func Fatalf(msg string, args ...any) {
	l.logger.Sugar().Fatalf(msg, args)
}

func Fatalw(msg string, keyValues ...any) {
	l.logger.Sugar().Fatalw(msg, keyValues)
}

func DPanic(args ...any) {
	l.logger.Sugar().DPanic(args)
}

func DPanicf(msg string, args ...any) {
	l.logger.Sugar().DPanicf(msg, args)
}

func DPanicw(msg string, keyValues ...any) {
	l.logger.Sugar().DPanicw(msg, keyValues)
}

func Panic(args ...any) {
	l.logger.Sugar().Panic(args)
}

func Panicf(msg string, args ...any) {
	l.logger.Sugar().Panicf(msg, args)
}

func Panicw(msg string, keyValues ...any) {
	l.logger.Sugar().Panicw(msg, keyValues)
}

func ZapLogger() *zap.Logger {
	return l.logger
}
