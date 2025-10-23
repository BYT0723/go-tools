package osx

import (
	"os"
	"reflect"
	"strconv"

	"golang.org/x/exp/constraints"
)

func GetEnv[T constraints.Integer | constraints.Float | ~string | ~bool](key string, def T) T {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	t := reflect.TypeOf(def)
	vt := reflect.New(t).Elem()

	switch t.Kind() {
	case reflect.String:
		vt.SetString(v)
	case reflect.Bool:
		val, err := strconv.ParseBool(v)
		if err != nil {
			return def
		}
		vt.SetBool(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return def
		}
		vt.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return def
		}
		vt.SetUint(val)
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return def
		}
		vt.SetFloat(val)
	}
	return vt.Interface().(T)
}
