package logx

type (
	LoggerType uint8
)

const (
	TypeZap LoggerType = iota // default
	TypeZeroLog

	TypeInvalid LoggerType = 1<<4 - 1
)
