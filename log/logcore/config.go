package logcore

type (
	LoggerType uint8
)

const (
	ZAP = iota // default
	ZEROLOG
	INVALID
)

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
	AllIn bool
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

func DefaultLoggerConf() *LoggerConf {
	return &LoggerConf{
		Dir:        "logs",
		Name:       "app",
		Ext:        "log",
		Level:      "debug",
		AllIn:      false,
		MaxBackups: 20,
		MaxSize:    20,
		MaxAge:     7,
		Console:    true,
	}
}

func (cfg *LoggerConf) Clone() *LoggerConf {
	res := *cfg
	return &res
}
