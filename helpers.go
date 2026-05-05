package defaults

import (
	"errors"
)

// Value initializes a Provider[T] with a fallback value.
//
// Example:
//
//	ageDef := Value(21)
func Value[T any](val T) Provider[T] {
	return Provider[T]{defaultValue: val}
}

// Normalize ensures a slice has a minimum length by padding it with nil.
// This prevents "index out of range" panics when processing variadic arguments.
//
// Example:
//
//	args := Normalize([]any{1}, 3)
//	// args is now []any{1, nil, nil}
func Normalize(values []any, needed int) []any {
	neededArgs := make([]any, needed)
	copy(neededArgs, values)
	return neededArgs
}

// Aggregate collects multiple results and joins them into a single error.
// Returns nil if all status.Ok flags are true.
//
// Example:
//
//	err := Aggregate(status1, status2)
func Aggregate(args ...Result) error {
	errList := make([]error, 0, len(args))

	for _, arg := range args {
		if !arg.Ok {
			errList = append(errList, arg)
		}
	}

	if len(errList) > 0 {
		return errors.Join(errList...)
	}

	return nil
}
