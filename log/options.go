package log

type (
	LoggerType uint8
	InitOption func(cfg *InitConf)
)

const (
	ZEROLOG = iota
	ZAP
)

type InitConf struct {
	Type LoggerType
}

func WithLoggerType(_t LoggerType) InitOption {
	return func(cfg *InitConf) {
		cfg.Type = _t
	}
}
