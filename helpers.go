package defaults

import (
	"errors"
	"reflect"
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
	if len(values) >= needed {
		return values
	}
	neededArgs := make([]any, needed)
	copy(neededArgs, values)
	return neededArgs
}

// AggregateErrors collects multiple results and joins them into a single error.
// Returns nil if all status.Ok flags are true.
//
// Example:
//
//	err := AggregateErrors(status1, status2)
func AggregateErrors(args ...Result) error {
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

// isTypedNil provides a fast-path for detecting nil pointers wrapped in interfaces.
// It handles common nilable types via type switch and falls back to reflection for others.
func isTypedNil(value any) bool {
	if value == nil {
		return true
	}
	switch v := value.(type) {
	case *int:
		return v == nil
	case *int8:
		return v == nil
	case *int16:
		return v == nil
	case *int32:
		return v == nil
	case *int64:
		return v == nil
	case *uint:
		return v == nil
	case *uint8:
		return v == nil
	case *uint16:
		return v == nil
	case *uint32:
		return v == nil
	case *uint64:
		return v == nil
	case *uintptr:
		return v == nil
	case *float32:
		return v == nil
	case *float64:
		return v == nil
	case *bool:
		return v == nil
	case *string:
		return v == nil
	case *complex64:
		return v == nil
	case *complex128:
		return v == nil
	case []any:
		return v == nil
	case []int:
		return v == nil
	case []string:
		return v == nil
	case []bool:
		return v == nil
	case map[string]any:
		return v == nil
	case map[string]int:
		return v == nil
	case map[int]string:
		return v == nil
	case chan any:
		return v == nil
	case chan int:
		return v == nil
	case chan string:
		return v == nil
	case func():
		return v == nil
	case func(int) int:
		return v == nil
	default:
		// Fallback to reflection for uncovered types
		rv := reflect.ValueOf(value)
		if rv.IsValid() {
			switch rv.Kind() {
			case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
				return rv.IsNil()
			}
		}
		return false
	}
}
