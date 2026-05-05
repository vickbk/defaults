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
	typeStatus := DefaultType{Ok: true, UsedDefault: true}

	if input == nil {
		return d.defaultValue, typeStatus
	}

	if val, ok := input.(T); ok {
		typeStatus.UsedDefault = false
		return val, typeStatus
	}

	// Type mismatch: Return default and 'false' for correctness
	v := reflect.ValueOf(input)
	k := v.Kind()

	// Only call IsNil on types that support it to avoid panics
	if k == reflect.Ptr || k == reflect.Map || k == reflect.Slice || k == reflect.Chan || k == reflect.Interface || k == reflect.Func {
		if v.IsNil() {
			return d.defaultValue, typeStatus
		}
	}

	typeStatus.Ok = false
	typeStatus.Message = Optional(message, fmt.Sprintf("Invalid type for %T", d.defaultValue))

	return d.defaultValue, typeStatus
}

func (d Default[T]) SafeCheck(values []any, index int, message ...string) (T, DefaultType) {
	value := any(nil)
	if index < len(values) {
		value = values[index]
	}

	return d.Check(value, message...)
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
