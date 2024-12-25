package log

import (
	"testing"

	"github.com/BYT0723/go-tools/log/logcore"
)

func BenchmarkZap(b *testing.B) {
	l, err := NewLogger(WithLoggerType(
		TypeZap),
		WithConf(&logcore.LoggerConf{
			Dir:        "logs",
			Name:       "app",
			Ext:        ".log",
			Level:      "debug",
			Single:     false,
			MaxBackups: 5,
			MaxSize:    10,
			MaxAge:     7,
			Console:    false,
		}),
	)
	if err != nil {
		b.Fail()
	}
	for i := 0; i < b.N; i++ {
		l.Debug(i)
	}
}

func BenchmarkZeroLog(b *testing.B) {
	l, err := NewLogger(
		WithLoggerType(TypeZeroLog),
		WithConf(&logcore.LoggerConf{
			Dir:        "logs",
			Name:       "app",
			Ext:        ".log",
			Level:      "debug",
			Single:     false,
			MaxBackups: 5,
			MaxSize:    10,
			MaxAge:     7,
			Console:    false,
		}),
	)
	if err != nil {
		b.Fail()
	}
	for i := 0; i < b.N; i++ {
		l.Debug(i)
	}
}
