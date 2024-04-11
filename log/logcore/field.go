package logcore

type Field struct {
	Key   string
	Value any
}

func Any(key string, value any) *Field {
	return &Field{
		Key:   key,
		Value: value,
	}
}
