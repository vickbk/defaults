package defaults

// Optional returns the first element of a slice or a defaultValue if the slice is empty.
//
// Example:
//
//	msg := Optional(customMessages, "Standard Error")
func Optional[T any](values []T, defaultValue T) T {
	if len(values) > 0 {
		return values[0]
	}

	return defaultValue
}
// Optionals returns a slice of values where each element is taken from the input slice if it exists, or from the defaultValues if the input slice does not have enough elements.
//
// Example:
//
//	msgs := Optionals(customMessages, "Error 1", "Error 2", "Error 3")
func Optionals[T any](values []T, defaultValues ...T) []T {
	results := make([]T, len(defaultValues))

	valueLen := len(values)

	for i, defaultValue := range defaultValues {
		if i < valueLen {
			results[i] = values[i]
		} else {
			results[i] = defaultValue
		}
	}

	return results
}
