package zerologger

import (
	"fmt"
	"io"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/BYT0723/go-tools/logx/logcore"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type zeroLogger struct {
	zero zerolog.Logger
}

var defaultCallerSkip = 5

func NewInstance(cfg *logcore.LoggerConf) (ins *zeroLogger, err error) {
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFieldName = "timestamp"
	zerolog.MessageFieldName = "msg"
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		defaultPath := file + ":" + strconv.Itoa(line)
		idx := strings.LastIndexByte(file, '/')
		if idx == -1 {
			return defaultPath
		}
		idx = strings.LastIndexByte(file[:idx], '/')
		if idx == -1 {
			return defaultPath
		}
		return file[idx+1:] + ":" + strconv.Itoa(line)
	}

	var (
		writers  []io.Writer
		basename = filepath.Join(cfg.Dir, cfg.Name)
	)

	if !cfg.Multi {
		writers = append(writers, NewLevelWriter(zerolog.SyncWriter(&lumberjack.Logger{
			Filename:   basename + cfg.Ext,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
		}), func(l zerolog.Level) bool { return l >= level }))
	} else {
		for i := level; i < zerolog.Disabled; i++ {
			var (
				targetLevel = i
				filename    = basename + "-" + targetLevel.String() + cfg.Ext
			)
			writers = append(writers, NewLevelWriter(zerolog.SyncWriter(&lumberjack.Logger{
				Filename:   filename,
				MaxSize:    cfg.MaxSize,
				MaxBackups: cfg.MaxBackups,
				MaxAge:     cfg.MaxAge,
			}), func(l zerolog.Level) bool { return l == targetLevel }))
		}
	}

	if cfg.Console {
		writers = append(writers, zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = zerolog.TimeFieldFormat + "\t"
			w.FormatLevel = func(i any) string {
				return strings.ToUpper(fmt.Sprint(i, "\t"))
			}
			w.FormatMessage = func(i any) string {
				return fmt.Sprint(i, "\t")
			}
		}))
	}

	ins = &zeroLogger{
		zero: zerolog.New(zerolog.MultiLevelWriter(writers...)).
			With().
			Timestamp().
			CallerWithSkipFrameCount(defaultCallerSkip).
			Logger(),
	}

	return
}

func (l *zeroLogger) With(kvs ...logcore.Field) logcore.Logger {
	copy := l.clone()
	ctx := copy.zero.With()
	for _, kv := range kvs {
		switch kv.Kind {
		case reflect.Bool:
			ctx = ctx.Bool(kv.Key, kv.Value.(bool))
		case reflect.Int:
			ctx = ctx.Int(kv.Key, kv.Value.(int))
		case reflect.Int8:
			ctx = ctx.Int8(kv.Key, kv.Value.(int8))
		case reflect.Int16:
			ctx = ctx.Int16(kv.Key, kv.Value.(int16))
		case reflect.Int32:
			ctx = ctx.Int32(kv.Key, kv.Value.(int32))
		case reflect.Int64:
			ctx = ctx.Int64(kv.Key, kv.Value.(int64))
		case reflect.Uint:
			ctx = ctx.Uint(kv.Key, kv.Value.(uint))
		case reflect.Uint8:
			ctx = ctx.Uint8(kv.Key, kv.Value.(uint8))
		case reflect.Uint16:
			ctx = ctx.Uint16(kv.Key, kv.Value.(uint16))
		case reflect.Uint32:
			ctx = ctx.Uint32(kv.Key, kv.Value.(uint32))
		case reflect.Uint64:
			ctx = ctx.Uint64(kv.Key, kv.Value.(uint64))
		case reflect.Float32:
			ctx = ctx.Float32(kv.Key, kv.Value.(float32))
		case reflect.Float64:
			ctx = ctx.Float64(kv.Key, kv.Value.(float64))
		case reflect.String:
			ctx = ctx.Str(kv.Key, kv.Value.(string))
		default:
			switch v := kv.Value.(type) {
			case error:
				ctx = ctx.AnErr(kv.Key, v)
			case fmt.Stringer:
				ctx = ctx.Stringer(kv.Key, v)
			default:
				ctx = ctx.Any(kv.Key, kv.Value)
			}
		}
	}
	copy.zero = ctx.CallerWithSkipFrameCount(defaultCallerSkip - 1).Logger()
	return copy
}

func (l *zeroLogger) Debug(msg string, kvs ...logcore.Field) {
	l.log(zerolog.DebugLevel, msg, kvs...)
}

func (l *zeroLogger) Debugf(format string, args ...any) {
	l.logf(zerolog.DebugLevel, format, args...)
}

func (l *zeroLogger) Info(msg string, kvs ...logcore.Field) {
	l.log(zerolog.InfoLevel, msg, kvs...)
}

func (l *zeroLogger) Infof(format string, args ...any) {
	l.logf(zerolog.InfoLevel, format, args...)
}

func (l *zeroLogger) Warn(msg string, kvs ...logcore.Field) {
	l.log(zerolog.WarnLevel, msg, kvs...)
}

func (l *zeroLogger) Warnf(format string, args ...any) {
	l.logf(zerolog.WarnLevel, format, args...)
}

func (l *zeroLogger) Error(msg string, kvs ...logcore.Field) {
	l.log(zerolog.ErrorLevel, msg, kvs...)
}

func (l *zeroLogger) Errorf(format string, args ...any) {
	l.logf(zerolog.ErrorLevel, format, args...)
}

func (l *zeroLogger) Panic(msg string, kvs ...logcore.Field) {
	l.log(zerolog.PanicLevel, msg, kvs...)
}

func (l *zeroLogger) Panicf(format string, args ...any) {
	l.logf(zerolog.PanicLevel, format, args...)
}

func (l *zeroLogger) Fatal(msg string, kvs ...logcore.Field) {
	l.log(zerolog.FatalLevel, msg, kvs...)
}

func (l *zeroLogger) Fatalf(format string, args ...any) {
	l.logf(zerolog.FatalLevel, format, args...)
}

func (l *zeroLogger) Log(level string, msg string, kvs ...logcore.Field) {
	lv, err := zerolog.ParseLevel(level)
	if err != nil {
		lv = zerolog.DebugLevel
	}
	l.log(lv, msg, kvs...)
}

func (l *zeroLogger) Logf(level, format string, args ...any) {
	lv, err := zerolog.ParseLevel(level)
	if err != nil {
		lv = zerolog.DebugLevel
	}
	l.zero.WithLevel(lv).Msgf(format, args...)
}

func (l *zeroLogger) Sync() error {
	return nil
}

func (l *zeroLogger) Logger() logcore.Logger {
	l.zero = l.zero.With().CallerWithSkipFrameCount(defaultCallerSkip - 1).Logger()
	return l
}

func (l *zeroLogger) clone() *zeroLogger {
	copy := *l
	return &copy
}

func (l *zeroLogger) AddCallerSkip(skip int) logcore.Logger {
	zl := l.clone()
	zl.zero = zl.zero.With().CallerWithSkipFrameCount(defaultCallerSkip + skip).Logger()
	return zl
}

func (l *zeroLogger) log(lv zerolog.Level, msg string, kvs ...logcore.Field) {
	e := l.zero.WithLevel(lv)
	addFields(e, kvs...)
	e.Msg(msg)
}

func (l *zeroLogger) logf(lv zerolog.Level, format string, args ...any) {
	l.zero.WithLevel(lv).Msgf(format, args...)
}

func addFields(e *zerolog.Event, kvs ...logcore.Field) {
	for _, kv := range kvs {
		switch kv.Kind {
		case reflect.Bool:
			e.Bool(kv.Key, kv.Value.(bool))
		case reflect.Int:
			e.Int(kv.Key, kv.Value.(int))
		case reflect.Int8:
			e.Int8(kv.Key, kv.Value.(int8))
		case reflect.Int16:
			e.Int16(kv.Key, kv.Value.(int16))
		case reflect.Int32:
			e.Int32(kv.Key, kv.Value.(int32))
		case reflect.Int64:
			e.Int64(kv.Key, kv.Value.(int64))
		case reflect.Uint:
			e.Uint(kv.Key, kv.Value.(uint))
		case reflect.Uint8:
			e.Uint8(kv.Key, kv.Value.(uint8))
		case reflect.Uint16:
			e.Uint16(kv.Key, kv.Value.(uint16))
		case reflect.Uint32:
			e.Uint32(kv.Key, kv.Value.(uint32))
		case reflect.Uint64:
			e.Uint64(kv.Key, kv.Value.(uint64))
		case reflect.Float32:
			e.Float32(kv.Key, kv.Value.(float32))
		case reflect.Float64:
			e.Float64(kv.Key, kv.Value.(float64))
		case reflect.String:
			e.Str(kv.Key, kv.Value.(string))
		default:
			switch v := kv.Value.(type) {
			case error:
				e.AnErr(kv.Key, v)
			case fmt.Stringer:
				e.Stringer(kv.Key, v)
			default:
				e.Any(kv.Key, kv.Value)
			}
		}
	}
}
