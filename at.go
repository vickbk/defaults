package defaults

// At returns the element at the specified index from the input slice or a defaultValue
// if the index is out of bounds. Does not panic if index is negative (bounds check).
//
// This is the fastest path with zero allocations and no reflection overhead.
//
// Example:
//
//	modes := []string{"debug", "info", "error"}
//	mode := defaults.At(modes, 1, "unknown") // Returns "info"
//	mode := defaults.At(modes, 10, "unknown") // Returns "unknown"
func At[T any](values []T, index int, defaultValue T) T {
	if index >= 0 && index < len(values) {
		return values[index]
	}

	return defaultValue
}
