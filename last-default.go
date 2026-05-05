package defaults

func Optional[T any](values []T, defaultValue T) T {
	if len(values) > 0 {
		return values[0]
	}

	return defaultValue
}
