package defaults

// OptionalAt returns the element at the specified index from the input slice or a defaultValue if the index is out of bounds.
//
// Example:
//
//	msg := OptionalAt(customMessages, 0, "Standard Error")
//
// Deprecated: Use At instead for a classic naming
func OptionalAt[T any](values []T, index int, defaultValue T) T {
	return At(values, index, defaultValue)
}

// Optional returns the first element of a slice or a defaultValue if the slice is empty.
//
// Example:
//
//	msg := Optional(customMessages, "Standard Error")
//
// Deprecated: use Get instead
func Optional[T any](values []T, defaultValue T) T {
	return Get(values, defaultValue)
}

// Optionals returns a slice of values where each element is taken from the input slice if it exists, or from the defaultValues if the input slice does not have enough elements.
//
// Example:
//
//	msgs := Optionals(customMessages, "Error 1", "Error 2", "Error 3")
//
// Deprecated: use Slice instead
func Optionals[T any](values []T, defaultValues ...T) []T {

	return Slice(values, defaultValues...)
}
