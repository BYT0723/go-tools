package zerologger

import (
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"github.com/BYT0723/go-tools/log/logger"
	"github.com/rs/zerolog"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

type zeroLogger struct {
	logger *zerolog.Logger
	cfg    *logger.LoggerConf
}

func NewInstance(opts ...logger.Option) (ins *zeroLogger, err error) {
	ins = &zeroLogger{cfg: logger.DefaultLoggerConf()}

	for _, opt := range opts {
		opt(ins.cfg)
	}

	level, err := zerolog.ParseLevel(ins.cfg.Level)
	if err != nil {
		return nil, err
	}

	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFieldName = "timestamp"

	var writers []io.Writer

	if ins.cfg.AllIn {
		writers = append(writers, NewLevelWriter(zerolog.SyncWriter(&lumberjack.Logger{
			Filename:   path.Join(ins.cfg.Dir, fmt.Sprintf("%s.%s", ins.cfg.Name, ins.cfg.Ext)),
			MaxSize:    100,
			MaxBackups: 5,
			MaxAge:     7,
		}), func(l zerolog.Level) bool { return l >= level }))
	} else {
		for i := level; i < zerolog.Disabled; i++ {
			targetLevel := i
			writers = append(writers, NewLevelWriter(zerolog.SyncWriter(&lumberjack.Logger{
				Filename:   path.Join(ins.cfg.Dir, fmt.Sprintf("%s-%s.%s", ins.cfg.Name, targetLevel, ins.cfg.Ext)),
				MaxSize:    100,
				MaxBackups: 5,
				MaxAge:     7,
			}), func(l zerolog.Level) bool { return l == targetLevel }))
		}
	}

	if ins.cfg.Console {
		writers = append(writers, zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = zerolog.TimeFieldFormat
			w.FormatLevel = func(i interface{}) string {
				return strings.ToUpper(fmt.Sprint(i) + "\t")
			}
		}))
	}

	l := zerolog.New(zerolog.MultiLevelWriter(writers...)).With().Timestamp().Caller().Logger()
	ins.logger = &l

	return
}

func (l *zeroLogger) With(kvs ...*logger.Field) logger.Logger {
	copy := l.Clone()
	ctx := copy.logger.With()
	for _, v := range kvs {
		ctx = ctx.Any(v.Key, v.Value)
	}
	nlog := ctx.Logger()
	copy.logger = &nlog
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

func (l *zeroLogger) Sync() error {
	return nil
}

func (l *zeroLogger) Clone() *zeroLogger {
	copy := *l
	cpLogger := *copy.logger
	copy.logger = &cpLogger
	return &copy
}
