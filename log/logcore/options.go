package logcore

type Option func(cfg *InitConf)

type InitConf struct {
	Type   LoggerType
	LogCfg *LoggerConf
}

func WithLoggerType(_t LoggerType) Option {
	return func(cfg *InitConf) {
		if _t == 0 || _t >= INVALID {
			_t = ZAP
		}
		cfg.Type = _t
	}
}

func WithLevel(level string) Option {
	return func(cfg *InitConf) {
		cfg.LogCfg.Level = level
	}
}

func WithName(name string) Option {
	return func(cfg *InitConf) {
		cfg.LogCfg.Name = name
	}
}

func WithPath(dir string) Option {
	return func(cfg *InitConf) {
		cfg.LogCfg.Dir = dir
	}
}

func WithMaxBackups(backups int) Option {
	return func(cfg *InitConf) {
		cfg.LogCfg.MaxBackups = backups
	}
}

// size: mb
func WithMaxSize(size int) Option {
	return func(cfg *InitConf) {
		cfg.LogCfg.MaxSize = size
	}
}

// age: day
func WithMaxAge(age int) Option {
	return func(cfg *InitConf) {
		cfg.LogCfg.MaxAge = age
	}
}

func WithConf(cfg *LoggerConf) Option {
	return func(payload *InitConf) {
		payload.LogCfg = cfg
	}
}
