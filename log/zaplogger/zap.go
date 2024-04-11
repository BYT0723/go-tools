package zaplogger

import (
	"fmt"
	"path/filepath"

	"github.com/BYT0723/go-tools/log/logcore"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *zap.Logger
}

func NewInstance(cfg *logcore.LoggerConf) (ins *zapLogger, err error) {
	level, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	cores := []zapcore.Core{}
	if cfg.AllIn {
		targetLevel := level.Level()
		cores = append(cores, newCore(
			cfg,
			func(l zapcore.Level) bool { return l >= targetLevel },
			filepath.Join(cfg.Dir, fmt.Sprintf("%s.%s", cfg.Name, cfg.Ext)),
		))
	} else {
		for i := level.Level(); i <= zap.FatalLevel; i++ {
			targetLevel := i
			cores = append(cores, newCore(
				cfg,
				func(l zapcore.Level) bool { return l == targetLevel },
				filepath.Join(cfg.Dir, fmt.Sprintf("%s-%s.%s", cfg.Name, targetLevel, cfg.Ext)),
			))
		}
	}

	if cfg.Console {
		cores = append(cores, newConsoleCore(level))
	}

	core := zapcore.NewTee(cores...)

	ins = &zapLogger{
		logger: zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)),
	}

	return
}

func (logger *zapLogger) With(kvs ...*logcore.Field) logcore.Logger {
	fields := []zap.Field{}
	for _, kv := range kvs {
		fields = append(fields, zap.Any(kv.Key, kv.Value))
	}

	res := logger.Clone()
	res.logger = res.logger.With(fields...)
	return res
}

func (logger *zapLogger) Debug(args ...any) {
	logger.logger.Sugar().Debug(args...)
}

func (logger *zapLogger) Debugf(format string, args ...any) {
	logger.logger.Sugar().Debugf(format, args...)
}

func (logger *zapLogger) Info(args ...any) {
	logger.logger.Sugar().Info(args...)
}

func (logger *zapLogger) Infof(format string, args ...any) {
	logger.logger.Sugar().Infof(format, args...)
}

func (logger *zapLogger) Warn(args ...any) {
	logger.logger.Sugar().Warn(args...)
}

func (logger *zapLogger) Warnf(format string, args ...any) {
	logger.logger.Sugar().Warnf(format, args)
}

func (logger *zapLogger) Error(args ...any) {
	logger.logger.Sugar().Error(args...)
}

func (logger *zapLogger) Errorf(format string, args ...any) {
	logger.logger.Sugar().Errorf(format, args...)
}

func (logger *zapLogger) Panic(args ...any) {
	logger.logger.Sugar().Panic(args...)
}

func (logger *zapLogger) Panicf(format string, args ...any) {
	logger.logger.Sugar().Panicf(format, args...)
}

func (logger *zapLogger) Fatal(args ...any) {
	logger.logger.Sugar().Fatal(args...)
}

func (logger *zapLogger) Fatalf(format string, args ...any) {
	logger.logger.Sugar().Fatalf(format, args...)
}

func (logger *zapLogger) ZapLogger() (*zap.Logger, bool) {
	var res bool
	if logger.logger != nil {
		res = true
	}
	return logger.logger, res
}

func (logger *zapLogger) Sync() error {
	return logger.logger.Sync()
}

func (logger *zapLogger) Clone() *zapLogger {
	copy := *logger
	cplog := *copy.logger
	copy.logger = &cplog
	return &copy
}
