package defaults

// Optional returns the first element of a slice or a defaultValue if the slice is empty.
//
// Example:
//   msg := Optional(customMessages, "Standard Error")
func Optional[T any](values []T, defaultValue T) T {
	if len(values) > 0 {
		return values[0]
	}

	return defaultValue
}
