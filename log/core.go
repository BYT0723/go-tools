package log

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func newConsoleCore(level zap.AtomicLevel) zapcore.Core {
	devCfg := zap.NewDevelopmentEncoderConfig()
	devCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(devCfg),
		zapcore.AddSync(os.Stdout),
		level,
	)
}

func newCore(level zap.AtomicLevel, onlyLevel bool, name string, cfg *LoggerConf) zapcore.Core {
	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339Nano)

	var levelEnabler zapcore.LevelEnabler = level
	if onlyLevel {
		levelEnabler = zap.LevelEnablerFunc(func(l zapcore.Level) bool { return l == level.Level() })
	}

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(productionCfg),
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   name,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
		}),
		levelEnabler,
	)
}
