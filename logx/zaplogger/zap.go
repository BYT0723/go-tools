package zaplogger

import (
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/BYT0723/go-tools/logx/logcore"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	zap   *zap.Logger
	sugar *zap.SugaredLogger
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

	zl := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))

	ins = &zapLogger{
		zap:   zl,
		sugar: zl.Sugar(),
	}

	return
}

func (l *zapLogger) With(kvs ...logcore.Field) logcore.Logger {
	zl := l.zap.With(transFields(kvs)...)
	return &zapLogger{
		zap:   zl,
		sugar: zl.Sugar(),
	}
}

func (l *zapLogger) Debug(msg string, kvs ...logcore.Field) {
	l.zap.Debug(msg, transFields(kvs)...)
}

func (l *zapLogger) Debugf(format string, args ...any) {
	l.sugar.Debugf(format, args...)
}

func (l *zapLogger) Info(msg string, kvs ...logcore.Field) {
	l.zap.Info(msg, transFields(kvs)...)
}

func (l *zapLogger) Infof(format string, args ...any) {
	l.sugar.Infof(format, args...)
}

func (l *zapLogger) Warn(msg string, kvs ...logcore.Field) {
	l.zap.Warn(msg, transFields(kvs)...)
}

func (l *zapLogger) Warnf(format string, args ...any) {
	l.sugar.Warnf(format, args)
}

func (l *zapLogger) Error(msg string, kvs ...logcore.Field) {
	l.zap.Error(msg, transFields(kvs)...)
}

func (l *zapLogger) Errorf(format string, args ...any) {
	l.sugar.Errorf(format, args...)
}

func (l *zapLogger) Panic(msg string, kvs ...logcore.Field) {
	l.zap.Panic(msg, transFields(kvs)...)
}

func (l *zapLogger) Panicf(format string, args ...any) {
	l.sugar.Panicf(format, args...)
}

func (l *zapLogger) Fatal(msg string, kvs ...logcore.Field) {
	l.zap.Fatal(msg, transFields(kvs)...)
}

func (l *zapLogger) Fatalf(format string, args ...any) {
	l.sugar.Fatalf(format, args...)
}

func (l *zapLogger) Log(level string, msg string, kvs ...logcore.Field) {
	var lv zapcore.Level
	if v, err := zap.ParseAtomicLevel(level); err != nil {
		lv = zap.DebugLevel
	} else {
		lv = v.Level()
	}
	if ce := l.zap.WithOptions(zap.AddCallerSkip(1)).Check(lv, msg); ce != nil {
		ce.Write(transFields(kvs)...)
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

func transFields(fields []logcore.Field) (result []zapcore.Field) {
	for _, kv := range fields {
		switch kv.Kind {
		case reflect.Bool:
			result = append(result, zap.Bool(kv.Key, kv.Value.(bool)))
		case reflect.Int:
			result = append(result, zap.Int(kv.Key, kv.Value.(int)))
		case reflect.Int8:
			result = append(result, zap.Int8(kv.Key, kv.Value.(int8)))
		case reflect.Int16:
			result = append(result, zap.Int16(kv.Key, kv.Value.(int16)))
		case reflect.Int32:
			result = append(result, zap.Int32(kv.Key, kv.Value.(int32)))
		case reflect.Int64:
			result = append(result, zap.Int64(kv.Key, kv.Value.(int64)))
		case reflect.Uint:
			result = append(result, zap.Uint(kv.Key, kv.Value.(uint)))
		case reflect.Uint8:
			result = append(result, zap.Uint8(kv.Key, kv.Value.(uint8)))
		case reflect.Uint16:
			result = append(result, zap.Uint16(kv.Key, kv.Value.(uint16)))
		case reflect.Uint32:
			result = append(result, zap.Uint32(kv.Key, kv.Value.(uint32)))
		case reflect.Uint64:
			result = append(result, zap.Uint64(kv.Key, kv.Value.(uint64)))
		case reflect.Float32:
			result = append(result, zap.Float32(kv.Key, kv.Value.(float32)))
		case reflect.Float64:
			result = append(result, zap.Float64(kv.Key, kv.Value.(float64)))
		case reflect.String:
			result = append(result, zap.String(kv.Key, kv.Value.(string)))
		default:
			switch v := kv.Value.(type) {
			case error:
				result = append(result, zap.NamedError(kv.Key, v))
			case fmt.Stringer:
				result = append(result, zap.Stringer(kv.Key, v))
			default:
				result = append(result, zap.Any(kv.Key, v))
			}
		}
	}
	return
}
