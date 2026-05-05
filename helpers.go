package defaults

import (
	"errors"
	"fmt"
	"reflect"
)

func DefaultArgsNormalize(values []any, needed int) []any {
	neededArgs := make([]any, needed)
	copy(neededArgs, values)
	return neededArgs
}

// Default initializes the provider with a fallback value.
func DefaultValue[T any](val T) Default[T] {
	return Default[T]{defaultValue: val}
}

// Check attempts to cast the input to type T.
// If it fails or is the zero value, it returns the default.
func (d Default[T]) Check(input any, message ...string) (T, DefaultType) {
	status := DefaultType{Ok: true, UsedDefault: true}

	if input == nil {
		return d.defaultValue, status
	}

	if val, ok := input.(T); ok {
		status.UsedDefault = false
		return val, status
	}

	// Type mismatch: Return default and 'false' for correctness
	v := reflect.ValueOf(input)

	if v.IsValid() {

		switch v.Kind() {
		case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
			if v.IsNil() {
				return d.defaultValue, status
			}
		}

	}

	status.Ok = false
	status.Message = Optional(message, fmt.Sprintf("invalid type: expected %T, got %T", d.defaultValue, input))

	return d.defaultValue, status
}

func (d Default[T]) SafeCheck(values []any, index int, message ...string) (T, DefaultType) {
	value := any(nil)
	if index < len(values) {
		value = values[index]
	}

	return d.Check(value, message...)
}

func (d Default[T]) SafeCheckOrPanic(values []any, index int, message ...string) T {
	val, status := d.SafeCheck(values, index, message...)
	if !status.Ok {
		panic(status.Message)
	}
	return val
}

// CheckDefaults aggregates the correctness booleans.
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
