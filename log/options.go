package log

import (
	"regexp"
	"strings"

	"github.com/BYT0723/go-tools/log/logcore"
)

type Option func(cfg *InitConf)

type InitConf struct {
	Type   LoggerType
	LogCfg *logcore.LoggerConf
}

func WithLoggerType(_t LoggerType) Option {
	return func(cfg *InitConf) {
		if _t == 0 || _t >= TypeInvalid {
			_t = TypeZap
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

var extRegex = regexp.MustCompile(`^\.[a-zA-Z0-9]+$`)

func WithConf(cfg *logcore.LoggerConf) Option {
	return func(payload *InitConf) {
		cfg.Ext = extRegex.FindString(strings.TrimSpace(cfg.Ext))
		payload.LogCfg.Merge(cfg)
	}
}
