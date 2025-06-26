package osx

import (
	"os"
	"strconv"
)

func GetEnv[T any](key string, def T) T {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	var anyVal any
	switch any(def).(type) {
	case int:
		val, err := strconv.Atoi(v)
		if err != nil {
			return def
		}
		anyVal = val
	case int8:
		val, err := strconv.ParseInt(v, 10, 8)
		if err != nil {
			return def
		}
		anyVal = int8(val)
	case int16:
		val, err := strconv.ParseInt(v, 10, 16)
		if err != nil {
			return def
		}
		anyVal = int16(val)
	case int32:
		val, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return def
		}
		anyVal = int32(val)
	case int64:
		val, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return def
		}
		anyVal = val

	case uint:
		val, err := strconv.ParseUint(v, 10, 0)
		if err != nil {
			return def
		}
		anyVal = uint(val)
	case uint8:
		val, err := strconv.ParseUint(v, 10, 8)
		if err != nil {
			return def
		}
		anyVal = uint8(val)
	case uint16:
		val, err := strconv.ParseUint(v, 10, 16)
		if err != nil {
			return def
		}
		anyVal = uint16(val)
	case uint32:
		val, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return def
		}
		anyVal = uint32(val)
	case uint64:
		val, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return def
		}
		anyVal = val

	case float32:
		val, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return def
		}
		anyVal = float32(val)
	case float64:
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return def
		}
		anyVal = val

	case string:
		anyVal = v

	case bool:
		val, err := strconv.ParseBool(v)
		if err != nil {
			return def
		}
		anyVal = val

	default:
		return def
	}

	return anyVal.(T)
}
