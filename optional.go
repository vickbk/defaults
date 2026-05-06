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
	vLen := len(values)
	dLen := len(defaultValues)

	// If the user provided enough (or more) values, return the input slice.
	if vLen >= dLen {
		return values
	}

	// We allocate exactly dLen because we know that's the required size.
	results := make([]T, dLen)

	// Copy original values into the start of results
	copy(results, values)

	for i := vLen; i < dLen; i++ {
		results[i] = defaultValues[i]
	}

	return results
}

// OptionalAt returns the element at the specified index from the input slice or a defaultValue if the index is out of bounds.
//
// Example:
//
//	msg := OptionalAt(customMessages, 0, "Standard Error")
func OptionalAt[T any](values []T, index int, defaultValue T) T {
	if index >= 0 && index < len(values) {
		return values[index]
	}

	return defaultValue
}
