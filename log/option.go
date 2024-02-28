package log

type Option func(logger *Logger)

func WithLevel(level string) Option {
	return func(logger *Logger) {
		logger.Cfg.Level = level
	}
}

func WithFileName(name string) Option {
	return func(logger *Logger) {
		logger.Cfg.Filename = name
	}
}

func WithMaxBackups(backups int) Option {
	return func(logger *Logger) {
		logger.Cfg.MaxBackups = backups
	}
}

// size: mb
func WithMaxSize(size int) Option {
	return func(logger *Logger) {
		logger.Cfg.MaxSize = size
	}
}

// age: day
func WithMaxAge(age int) Option {
	return func(logger *Logger) {
		logger.Cfg.MaxAge = age
	}
}

func WithConf(cfg *LoggerConf) Option {
	return func(logger *Logger) {
		logger.Cfg = cfg
	}
}
