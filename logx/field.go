package logx

import "github.com/BYT0723/go-tools/logx/logcore"

type Field = logcore.Field

func Any(key string, value any) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}
