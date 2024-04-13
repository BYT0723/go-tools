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

func (cfg *LoggerConf) Clone() *LoggerConf {
	res := *cfg
	return &res
}
