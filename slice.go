package defaults

// Slice returns a slice of values where each element is taken from the input slice if it exists, or from the defaultValues if the input slice does not have enough elements.
//
// Example:
//
//	msgs := Slice(customMessages, "Error 1", "Error 2", "Error 3")
func Slice[T any](values []T, defaultValues ...T) []T {
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
