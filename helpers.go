package defaults

import (
	"errors"
)

// DefaultValue initializes a Default[T] provider with a fallback value.
//
// Example:
//
//	ageDef := DefaultValue(21)
func DefaultValue[T any](val T) Default[T] {
	return Default[T]{defaultValue: val}
}

// DefaultArgsNormalize ensures a slice has a minimum length by padding it with nil.
// This prevents "index out of range" panics when processing variadic arguments.
//
// Example:
//
//	args := DefaultArgsNormalize([]any{1}, 3)
//	// args is now []any{1, nil, nil}
func DefaultArgsNormalize(values []any, needed int) []any {
	neededArgs := make([]any, needed)
	copy(neededArgs, values)
	return neededArgs
}

// CheckDefaults collects multiple DefaultType results and joins them into a single error.
// Returns nil if all status.Ok flags are true.
//
// Example:
//
//	err := CheckDefaults(status1, status2)
func CheckDefaults(args ...DefaultType) error {
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
