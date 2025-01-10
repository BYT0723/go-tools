package logcore

type LoggerConf struct {
	// log folder.
	// default: logs
	Dir string
	// log file name.
	// default: app
	Name string
	// log file ext.
	// default: log
	Ext string
	// The lowest level of log.
	// default: debug
	Level string
	// Whether different files will be stored according to the log level.
	// if true, all level logs will be stored in a single file. Otherwise store hierarchically.
	// default: false
	Single bool
	// The maximum number of backups.
	// default: 20
	MaxBackups int
	// The maximum size of a single file, in MB.
	// default: 20
	MaxSize int
	// The maximum storage duration of the file, in days.
	// default: 7
	MaxAge int
	// Whether the console outputs.
	// default: true
	Console bool
}

// 合并LoggerConf
func (c *LoggerConf) Merge(cfg *LoggerConf) {
	if cfg.Dir != "" {
		c.Dir = cfg.Dir
	}
	if cfg.Name != "" {
		c.Name = cfg.Name
	}
	if cfg.Ext != "" {
		c.Ext = cfg.Ext
	}
	if cfg.Level != "" {
		c.Level = cfg.Level
	}
	if cfg.Single {
		c.Single = cfg.Single
	}
	if cfg.MaxBackups != 0 {
		c.MaxBackups = cfg.MaxBackups
	}
	if cfg.MaxSize != 0 {
		c.MaxSize = cfg.MaxSize
	}
	if cfg.MaxAge != 0 {
		c.MaxAge = cfg.MaxAge
	}
	if cfg.Console {
		c.Console = cfg.Console
	}
}

func DefaultLoggerConf() *LoggerConf {
	return &LoggerConf{
		Dir:        "logs",
		Name:       "app",
		Ext:        ".log",
		Level:      "debug",
		Single:     false,
		MaxBackups: 20,
		MaxSize:    20,
		MaxAge:     7,
		Console:    true,
	}
}
