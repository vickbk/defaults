package defaults

// Get returns the first element of a slice or a defaultValue if the slice is empty.
// This is the fastest path with zero allocations and no reflection overhead.
//
// Example:
//
//	ports := []int{8080, 8443}
//	port := defaults.Get(ports, 3000) // Returns 8080
//
//	empty := []int{}
//	port := defaults.Get(empty, 3000) // Returns 3000
func Get[T any](values []T, defaultValue T) T {
	return At(values, 0, defaultValue)
}
