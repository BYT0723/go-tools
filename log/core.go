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

func newCore(level zap.AtomicLevel, cfg *LoggerConf) zapcore.Core {
	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339Nano)

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(productionCfg),
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
		}),
		level,
	)
}
