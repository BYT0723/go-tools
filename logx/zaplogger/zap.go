package zaplogger

import (
	"fmt"
	"path/filepath"

	"github.com/BYT0723/go-tools/logx/logcore"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	zap *zap.Logger
}

func NewInstance(cfg *logcore.LoggerConf) (ins *zapLogger, err error) {
	level, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	var (
		cores    = []zapcore.Core{}
		basename = filepath.Join(cfg.Dir, cfg.Name)
	)
	if !cfg.Multi {
		targetLevel := level.Level()
		cores = append(cores, newCore(
			cfg,
			func(l zapcore.Level) bool { return l >= targetLevel },
			basename+cfg.Ext,
		))
	} else {
		for i := level.Level(); i <= zap.FatalLevel; i++ {
			var (
				targetLevel = i
				filename    = basename + "-" + targetLevel.String() + cfg.Ext
			)
			cores = append(cores, newCore(
				cfg,
				func(l zapcore.Level) bool { return l == targetLevel },
				filename,
			))
		}
	}

	if cfg.Console {
		cores = append(cores, newConsoleCore(level))
	}

	core := zapcore.NewTee(cores...)

	ins = &zapLogger{
		zap: zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2)),
	}

	return
}

func (l *zapLogger) With(kvs ...logcore.Field) logcore.Logger {
	fields := []zap.Field{}
	for _, kv := range kvs {
		fields = append(fields, zap.Any(kv.Key, kv.Value))
	}
	return &zapLogger{
		zap: l.zap.With(fields...),
	}
}

func (l *zapLogger) Debug(args ...any) {
	l.zap.Sugar().Debug(args...)
}

func (l *zapLogger) Debugf(format string, args ...any) {
	l.zap.Sugar().Debugf(format, args...)
}

func (l *zapLogger) Info(args ...any) {
	l.zap.Sugar().Info(args...)
}

func (l *zapLogger) Infof(format string, args ...any) {
	l.zap.Sugar().Infof(format, args...)
}

func (l *zapLogger) Warn(args ...any) {
	l.zap.Sugar().Warn(args...)
}

func (l *zapLogger) Warnf(format string, args ...any) {
	l.zap.Sugar().Warnf(format, args)
}

func (l *zapLogger) Error(args ...any) {
	l.zap.Sugar().Error(args...)
}

func (l *zapLogger) Errorf(format string, args ...any) {
	l.zap.Sugar().Errorf(format, args...)
}

func (l *zapLogger) Panic(args ...any) {
	l.zap.Sugar().Panic(args...)
}

func (l *zapLogger) Panicf(format string, args ...any) {
	l.zap.Sugar().Panicf(format, args...)
}

func (l *zapLogger) Fatal(args ...any) {
	l.zap.Sugar().Fatal(args...)
}

func (l *zapLogger) Fatalf(format string, args ...any) {
	l.zap.Sugar().Fatalf(format, args...)
}

func (l *zapLogger) Log(level string, args ...any) {
	var lv zapcore.Level
	if v, err := zap.ParseAtomicLevel(level); err != nil {
		lv = zap.DebugLevel
	} else {
		lv = v.Level()
	}
	if ce := l.zap.WithOptions(zap.AddCallerSkip(1)).Check(lv, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func (l *zapLogger) Logf(level, format string, args ...any) {
	var lv zapcore.Level
	if v, err := zap.ParseAtomicLevel(level); err != nil {
		lv = zap.DebugLevel
	} else {
		lv = v.Level()
	}
	if ce := l.zap.Check(lv, fmt.Sprintf(format, args...)); ce != nil {
		ce.Write()
	}
}

func (l *zapLogger) Sync() error {
	return l.zap.Sync()
}

func (l *zapLogger) Logger() logcore.Logger {
	l.zap = l.zap.WithOptions(zap.AddCallerSkip(-1))
	return l
}

func (l *zapLogger) AddCallerSkip(skip int) logcore.Logger {
	zl := &zapLogger{zap: l.zap}
	zl.zap = zl.zap.WithOptions(zap.AddCallerSkip(skip))

	return zl
}
