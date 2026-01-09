package logx

import (
	"reflect"
	"time"

	"github.com/BYT0723/go-tools/logx/logcore"
)

type Field = logcore.Field

func Any(key string, value any) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}

func Bool(key string, value bool) Field {
	return Field{
		Key:   key,
		Kind:  reflect.Bool,
		Value: value,
	}
}

func String(key string, value string) Field {
	return Field{
		Key:   key,
		Kind:  reflect.String,
		Value: value,
	}
}

func Int(key string, value int) Field {
	return Field{
		Key:   key,
		Kind:  reflect.Int,
		Value: value,
	}
}

func Int8(key string, value int8) Field {
	return Field{
		Key:   key,
		Kind:  reflect.Int8,
		Value: value,
	}
}

func Int16(key string, value int16) Field {
	return Field{
		Key:   key,
		Kind:  reflect.Int16,
		Value: value,
	}
}

func Int32(key string, value int32) Field {
	return Field{
		Key:   key,
		Kind:  reflect.Int32,
		Value: value,
	}
}

func Int64(key string, value int64) Field {
	return Field{
		Key:   key,
		Kind:  reflect.Int64,
		Value: value,
	}
}

func Uint(key string, value uint) Field {
	return Field{
		Key:   key,
		Kind:  reflect.Uint,
		Value: value,
	}
}

func Uint8(key string, value uint8) Field {
	return Field{
		Key:   key,
		Kind:  reflect.Uint8,
		Value: value,
	}
}

func Uint16(key string, value uint16) Field {
	return Field{
		Key:   key,
		Kind:  reflect.Uint16,
		Value: value,
	}
}

func Uint32(key string, value uint32) Field {
	return Field{
		Key:   key,
		Kind:  reflect.Uint32,
		Value: value,
	}
}

func Uint64(key string, value uint64) Field {
	return Field{
		Key:   key,
		Kind:  reflect.Uint64,
		Value: value,
	}
}

func Float32(key string, value float32) Field {
	return Field{
		Key:   key,
		Kind:  reflect.Float32,
		Value: value,
	}
}

func Float64(key string, value float64) Field {
	return Field{
		Key:   key,
		Kind:  reflect.Float64,
		Value: value,
	}
}

func Err(value error) Field {
	return Field{
		Key:   "error",
		Value: value.Error(),
	}
}

func Duration(key string, value time.Duration) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}
