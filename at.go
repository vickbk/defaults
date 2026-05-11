package defaults

// At returns the element at the specified index from the input slice or a defaultValue if the index is out of bounds.
//
// Example:
//
//	msg := defaults.At(customMessages, 0, "Standard Error")
func At[T any](values []T, index int, defaultValue T) T {
	if index >= 0 && index < len(values) {
		return values[index]
	}

	return defaultValue
}
