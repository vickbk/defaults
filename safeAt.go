package defaults

import (
	"fmt"
)

// defaults.SafeAt performs an index-safe check on a slice.
// If the index is out of bounds, it automatically returns the default value of the T with the Results.
//
// Example:
//
//	val, status := defaults.SafeAt(vars,0,"Guest")
//
// It validates if an input matches type T.
// It returns the defaultValue and UsedDefault=true if the input is nil or a typed nil pointer.
// If the type is mismatched, it returns an error status with Ok=false.
// The optional message parameter allows for custom error messages on type mismatch.
//
// Example:
//
//	val, status := defaults.SafeAt([]int{}, 1, "not an int", "Age must be a number")
//
// Note: The function uses reflection to check for typed nil values, ensuring that it correctly identifies nil pointers and interfaces, which can be a common source of bugs in Go when dealing with slices of interfaces or pointers.
// For multiple optional parameters of different types, prefer using the defaults.Apply function with custom Applier functions for better type safety and error handling.
func SafeAt[T any](values []any, index int, defaultValue T, message ...string) (T, Result) {
	status := Result{Ok: true, UsedDefault: true}

	var value any
	if index >= 0 && index < len(values) {
		value = values[index]
	}

	if value == nil {
		return defaultValue, status
	}

	if val, ok := value.(T); ok {
		status.UsedDefault = false
		return val, status
	}

	// Check for typed nil before full reflection
	if isTypedNil(value) {
		return defaultValue, status
	}

	// Type mismatch: Return default and 'false' for correctness
	status.Ok = false
	status.Message = Get(message, fmt.Sprintf("invalid type: expected %T, got %T", defaultValue, value))

	return defaultValue, status
}
