package defaults

import (
	"fmt"
)

// SafeAt performs an index-safe check on a slice. If the index is out of bounds,
// it automatically returns the default value. It validates if an input matches type T.
// It returns the defaultValue and UsedDefault=true if the input is nil or a typed nil pointer.
// If the type is mismatched, it returns an error status with Ok=false.
//
// This function handles the Typed Nil Paradox using reflection fallback only when necessary.
// For strictly typed slices, prefer Get/At for better performance.
//
// Example:
//
//	val, status := SafeAt([]any{42}, 0, 0, "must be int")
//	// val=42, status.Ok=true
//
//	val, status := SafeAt([]any{"bad"}, 0, 0, "must be int")
//	// val=0, status.Ok=false, status.Message="invalid type: expected int, got string"
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
