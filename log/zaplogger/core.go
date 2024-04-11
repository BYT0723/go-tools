package zaplogger

import (
	"os"
	"time"

	"github.com/BYT0723/go-tools/log/logger"
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

func newCore(cfg *logger.LoggerConf, filter zap.LevelEnablerFunc, filename string) zapcore.Core {
	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339Nano)

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(productionCfg),
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
		}),
		filter,
	)
}
