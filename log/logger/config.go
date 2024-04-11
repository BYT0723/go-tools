package logger

type LoggerInitFunc func(opts ...Option) (Logger, error)

type LoggerConf struct {
	Dir        string
	Name       string
	Ext        string
	Level      string
	AllIn      bool
	MaxBackups int // uint: MB
	MaxSize    int
	MaxAge     int // uint: DAY
	Console    bool
}

func DefaultLoggerConf() *LoggerConf {
	return &LoggerConf{
		Dir:        "logs",
		Name:       "app",
		Ext:        "log",
		Level:      "debug",
		AllIn:      false,
		MaxBackups: 3,
		MaxSize:    10,
		MaxAge:     7,
		Console:    true,
	}
}

func (cfg *LoggerConf) Copy(src *LoggerConf) {
	cfg.Dir = src.Dir
	cfg.Name = src.Name
	cfg.Ext = src.Ext
	cfg.Level = src.Level
	cfg.AllIn = src.AllIn
	cfg.MaxBackups = src.MaxBackups
	cfg.MaxSize = src.MaxSize
	cfg.MaxAge = src.MaxAge
}
