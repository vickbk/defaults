package defaults

func GetLastDefaultValue[T any](values []T, defaultValue ...T) T {
	if len(values) > 0 {
		return values[0]
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	var zero T

	return zero
}
