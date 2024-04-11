package zaplogger

import (
	"fmt"
	"path/filepath"

	"github.com/BYT0723/go-tools/log/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	sugarLogger *zap.SugaredLogger
	cfg         *logger.LoggerConf
}

func NewInstance(opts ...logger.Option) (ins *zapLogger, err error) {
	ins = &zapLogger{cfg: logger.DefaultLoggerConf()}

	for _, opt := range opts {
		opt(ins.cfg)
	}

	level, err := zap.ParseAtomicLevel(ins.cfg.Level)
	if err != nil {
		return nil, err
	}

	cores := []zapcore.Core{}
	if ins.cfg.AllIn {
		targetLevel := level.Level()
		cores = append(cores, newCore(
			ins.cfg,
			func(l zapcore.Level) bool { return l >= targetLevel },
			filepath.Join(ins.cfg.Dir, fmt.Sprintf("%s.%s", ins.cfg.Name, ins.cfg.Ext)),
		))
	} else {
		for i := level.Level(); i <= zap.FatalLevel; i++ {
			targetLevel := i
			cores = append(cores, newCore(
				ins.cfg,
				func(l zapcore.Level) bool { return l == targetLevel },
				filepath.Join(ins.cfg.Dir, fmt.Sprintf("%s-%s.%s", ins.cfg.Name, i, ins.cfg.Ext)),
			))
		}
	}

	if ins.cfg.Console {
		cores = append(cores, newConsoleCore(level))
	}

	core := zapcore.NewTee(cores...)
	ins.sugarLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()

	return
}

func (logger *zapLogger) With(kvs ...*logger.Field) logger.Logger {
	fields := []zap.Field{}
	for _, kv := range kvs {
		fields = append(fields, zap.Any(kv.Key, kv.Value))
	}

	res := logger.Clone()
	res.sugarLogger = res.sugarLogger.Desugar().With(fields...).Sugar()
	return res
}

func (logger *zapLogger) Debug(args ...any) {
	logger.sugarLogger.Debug(args...)
}

func (logger *zapLogger) Debugf(format string, args ...any) {
	logger.sugarLogger.Debugf(format, args...)
}

func (logger *zapLogger) Info(args ...any) {
	logger.sugarLogger.Info(args...)
}

func (logger *zapLogger) Infof(format string, args ...any) {
	logger.sugarLogger.Infof(format, args...)
}

func (logger *zapLogger) Warn(args ...any) {
	logger.sugarLogger.Warn(args...)
}

func (logger *zapLogger) Warnf(format string, args ...any) {
	logger.sugarLogger.Warnf(format, args)
}

func (logger *zapLogger) Error(args ...any) {
	logger.sugarLogger.Error(args...)
}

func (logger *zapLogger) Errorf(format string, args ...any) {
	logger.sugarLogger.Errorf(format, args...)
}

func (logger *zapLogger) Panic(args ...any) {
	logger.sugarLogger.Panic(args...)
}

func (logger *zapLogger) Panicf(format string, args ...any) {
	logger.sugarLogger.Panicf(format, args...)
}

func (logger *zapLogger) Fatal(args ...any) {
	logger.sugarLogger.Fatal(args...)
}

func (logger *zapLogger) Fatalf(format string, args ...any) {
	logger.sugarLogger.Fatalf(format, args...)
}

func (logger *zapLogger) ZapLogger() (*zap.Logger, bool) {
	var res bool
	if logger.sugarLogger != nil {
		res = true
	}
	return logger.sugarLogger.Desugar(), res
}

func (logger *zapLogger) Sync() error {
	return logger.sugarLogger.Sync()
}

func (logger *zapLogger) Clone() *zapLogger {
	copy := *logger
	cpSugar := *copy.sugarLogger
	copy.sugarLogger = &cpSugar
	return &copy
}
