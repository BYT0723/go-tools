package zerologger

import (
	"fmt"
	"io"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/BYT0723/go-tools/log/logcore"
	"github.com/rs/zerolog"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

type zeroLogger struct {
	logger zerolog.Logger
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

	var writers []io.Writer

	if cfg.AllIn {
		writers = append(writers, NewLevelWriter(zerolog.SyncWriter(&lumberjack.Logger{
			Filename:   path.Join(cfg.Dir, fmt.Sprintf("%s.%s", cfg.Name, cfg.Ext)),
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
		}), func(l zerolog.Level) bool { return l >= level }))
	} else {
		for i := level; i < zerolog.Disabled; i++ {
			targetLevel := i
			writers = append(writers, NewLevelWriter(zerolog.SyncWriter(&lumberjack.Logger{
				Filename:   path.Join(cfg.Dir, fmt.Sprintf("%s-%s.%s", cfg.Name, targetLevel, cfg.Ext)),
				MaxSize:    cfg.MaxSize,
				MaxBackups: cfg.MaxBackups,
				MaxAge:     cfg.MaxAge,
			}), func(l zerolog.Level) bool { return l == targetLevel }))
		}
	}

	if cfg.Console {
		writers = append(writers, zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = zerolog.TimeFieldFormat + "\t"
			w.FormatLevel = func(i interface{}) string {
				return strings.ToUpper(fmt.Sprint(i, "\t"))
			}
			w.FormatMessage = func(i interface{}) string {
				return fmt.Sprint(i, "\t")
			}
		}))
	}

	ins = &zeroLogger{
		logger: zerolog.New(zerolog.MultiLevelWriter(writers...)).With().Timestamp().CallerWithSkipFrameCount(4).Logger(),
	}

	return
}

func (l *zeroLogger) With(kvs ...*logcore.Field) logcore.Logger {
	copy := l.Clone()
	ctx := copy.logger.With()
	for _, v := range kvs {
		ctx = ctx.Any(v.Key, v.Value)
	}
	copy.logger = ctx.Logger()
	return copy
}

func (l *zeroLogger) Debug(args ...any) {
	l.logger.Debug().MsgFunc(func() string {
		return fmt.Sprint(args...)
	})
}

func (l *zeroLogger) Debugf(format string, args ...any) {
	l.logger.Debug().Msgf(format, args...)
}

func (l *zeroLogger) Info(args ...any) {
	l.logger.Info().MsgFunc(func() string {
		return fmt.Sprint(args...)
	})
}

func (l *zeroLogger) Infof(format string, args ...any) {
	l.logger.Info().Msgf(format, args...)
}

func (l *zeroLogger) Warn(args ...any) {
	l.logger.Warn().MsgFunc(func() string {
		return fmt.Sprint(args...)
	})
}

func (l *zeroLogger) Warnf(format string, args ...any) {
	l.logger.Warn().Msgf(format, args...)
}

func (l *zeroLogger) Error(args ...any) {
	l.logger.Error().MsgFunc(func() string {
		return fmt.Sprint(args...)
	})
}

func (l *zeroLogger) Errorf(format string, args ...any) {
	l.logger.Error().Msgf(format, args...)
}

func (l *zeroLogger) Panic(args ...any) {
	l.logger.Panic().MsgFunc(func() string {
		return fmt.Sprint(args...)
	})
}

func (l *zeroLogger) Panicf(format string, args ...any) {
	l.logger.Panic().Msgf(format, args...)
}

func (l *zeroLogger) Fatal(args ...any) {
	l.logger.Fatal().MsgFunc(func() string {
		return fmt.Sprint(args...)
	})
}

func (l *zeroLogger) Fatalf(format string, args ...any) {
	l.logger.Fatal().Msgf(format, args...)
}

func (l *zeroLogger) ZapLogger() (*zap.Logger, bool) {
	return nil, false
}

func (l *zeroLogger) ZeroLogger() (*zerolog.Logger, bool) {
	l2 := l.logger.With().CallerWithSkipFrameCount(2).Logger()
	return &l2, true
}

func (l *zeroLogger) Sync() error {
	return nil
}

func (l *zeroLogger) Clone() *zeroLogger {
	copy := *l
	return &copy
}
