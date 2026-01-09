package logcore

import "reflect"

type Field struct {
	Key   string
	Kind  reflect.Kind
	Value any
}
