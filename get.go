package defaults

// Get returns the first element of a slice or a defaultValue if the slice is empty.
//
// Example:
//
//	msg := defaults.Get(customMessages, "Standard Error")
func Get[T any](values []T, defaultValue T) T {
	return At(values, 0, defaultValue)
}
