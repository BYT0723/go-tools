package logger

type Option func(logger *LoggerConf)

func WithLevel(level string) Option {
	return func(cfg *LoggerConf) {
		cfg.Level = level
	}
}

func WithName(name string) Option {
	return func(cfg *LoggerConf) {
		cfg.Name = name
	}
}

func WithPath(dir string) Option {
	return func(cfg *LoggerConf) {
		cfg.Dir = dir
	}
}

func WithMaxBackups(backups int) Option {
	return func(cfg *LoggerConf) {
		cfg.MaxBackups = backups
	}
}

// size: mb
func WithMaxSize(size int) Option {
	return func(cfg *LoggerConf) {
		cfg.MaxSize = size
	}
}

// age: day
func WithMaxAge(age int) Option {
	return func(cfg *LoggerConf) {
		cfg.MaxAge = age
	}
}

func WithConf(cfg *LoggerConf) Option {
	return func(payload *LoggerConf) {
		payload.Copy(cfg)
	}
}
