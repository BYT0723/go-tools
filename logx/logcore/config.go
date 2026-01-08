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
	// Split into different files according to different log levels
	// default: false
	Multi bool
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
	// default: false
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
	if cfg.MaxBackups > 0 {
		c.MaxBackups = cfg.MaxBackups
	}
	if cfg.MaxSize > 0 {
		c.MaxSize = cfg.MaxSize
	}
	if cfg.MaxAge > 0 {
		c.MaxAge = cfg.MaxAge
	}
	c.Multi = cfg.Multi
	c.Console = cfg.Console
}

func DefaultLoggerConf() *LoggerConf {
	return &LoggerConf{
		Dir:        "logs",
		Name:       "app",
		Ext:        ".log",
		Level:      "debug",
		Multi:      false,
		MaxBackups: 20,
		MaxSize:    20,
		MaxAge:     7,
		Console:    false,
	}
}
