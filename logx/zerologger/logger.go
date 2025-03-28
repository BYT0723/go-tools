package zerologger

import (
	"fmt"
	"io"
	"path/filepath"
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

	if cfg.Single {
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
			CallerWithSkipFrameCount(4).
			Logger(),
	}

	return
}

func (l *zeroLogger) With(kvs ...logcore.Field) logcore.Logger {
	copy := l.clone()
	ctx := copy.zero.With()
	for _, v := range kvs {
		ctx = ctx.Any(v.Key, v.Value)
	}
	copy.zero = ctx.Logger()
	return copy
}

func (l *zeroLogger) Debug(args ...any) {
	l.zero.Debug().MsgFunc(func() string {
		return fmt.Sprint(args...)
	})
}

func (l *zeroLogger) Debugf(format string, args ...any) {
	l.zero.Debug().Msgf(format, args...)
}

func (l *zeroLogger) Info(args ...any) {
	l.zero.Info().MsgFunc(func() string {
		return fmt.Sprint(args...)
	})
}

func (l *zeroLogger) Infof(format string, args ...any) {
	l.zero.Info().Msgf(format, args...)
}

func (l *zeroLogger) Warn(args ...any) {
	l.zero.Warn().MsgFunc(func() string {
		return fmt.Sprint(args...)
	})
}

func (l *zeroLogger) Warnf(format string, args ...any) {
	l.zero.Warn().Msgf(format, args...)
}

func (l *zeroLogger) Error(args ...any) {
	l.zero.Error().MsgFunc(func() string {
		return fmt.Sprint(args...)
	})
}

func (l *zeroLogger) Errorf(format string, args ...any) {
	l.zero.Error().Msgf(format, args...)
}

func (l *zeroLogger) Panic(args ...any) {
	l.zero.Panic().MsgFunc(func() string {
		return fmt.Sprint(args...)
	})
}

func (l *zeroLogger) Panicf(format string, args ...any) {
	l.zero.Panic().Msgf(format, args...)
}

func (l *zeroLogger) Fatal(args ...any) {
	l.zero.Fatal().MsgFunc(func() string {
		return fmt.Sprint(args...)
	})
}

func (l *zeroLogger) Fatalf(format string, args ...any) {
	l.zero.Fatal().Msgf(format, args...)
}

func (l *zeroLogger) Log(level string, args ...any) {
	lv, err := zerolog.ParseLevel(level)
	if err != nil {
		lv = zerolog.DebugLevel
	}
	l.zero.WithLevel(lv).MsgFunc(func() string {
		return fmt.Sprint(args...)
	})
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
	l.zero = l.zero.With().CallerWithSkipFrameCount(3).Logger()
	return l
}

func (l *zeroLogger) clone() *zeroLogger {
	copy := *l
	return &copy
}
