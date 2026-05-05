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

func (d Default[T]) IsDefault(values []any, index int) bool {
	if index >= len(values) || values[index] == nil {
		return true
	}

	_, t := d.Check(values[index])
	return !t.Ok
}

// Check attempts to cast the input to type T.
// If it fails or is the zero value, it returns the default.
func (d Default[T]) Check(input any, message ...string) (T, DefaultType) {
	typeStatus := DefaultType{
		Message: Optional(message, fmt.Sprintf("Invalid type for %T", d.defaultValue)),
		Ok:      true}

	if input == nil {
		return d.defaultValue, typeStatus
	}

	val, ok := input.(T)
	if !ok {
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
		return d.defaultValue, typeStatus
	}

	return val, typeStatus
}

func (d Default[T]) SafeCheck(values []any, index int, message ...string) (T, DefaultType) {
	if index >= len(values) {
		return d.defaultValue, DefaultType{
			Message: Optional(message, fmt.Sprintf("Invalid type for %T", d.defaultValue)),
			Ok:      true}
	}
	return d.Check(values[index], message...)
}

// CheckDefaults aggregates the correctness booleans.
func CheckDefaults(args ...DefaultType) error {
	errList := make([]error, 0, len(args))

	for _, arg := range args {
		if !arg.Ok {
			errList = append(errList, errors.New(arg.Message))
		}
	}

	if len(errList) > 0 {
		return errors.Join(errList...)
	}

	return nil
}
