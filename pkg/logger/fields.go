package logger

func Error(err error) Field {
	return Field{Key: "error", Value: err}
}

func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

func Int32(key string, value int32) Field {
	return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func String(key string, value string) Field {
	return Field{Key: key, Value: value}
}

func Any(key string, value any) Field {
	return Field{Key: key, Value: value}
}
