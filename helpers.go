package defaults

import (
	"errors"
	"reflect"
	"strings"
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
func (d Default[T]) Check(input any) (T, DefaultType) {
	typeStatus := DefaultType{Message: "", Ok: true}
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

func (d Default[T]) SafeCheck(values []any, index int) (T, DefaultType) {
	if index >= len(values) {
		return d.defaultValue, DefaultType{Message: "Index out of bounds", Ok: false}
	}
	return d.Check(values[index])
}

// CheckDefaults aggregates the correctness booleans.
func CheckDefaults(args ...any) error {
	errList := make([]string, 0, len(args)/2)

	for i := 0; i < len(args); i += 2 {
		ok := args[i].(bool)
		errMsg := args[i+1].(string)
		if !ok {
			errList = append(errList, errMsg)
		}
	}

	if len(errList) > 0 {
		return errors.New(strings.Join(errList, ","))
	}

	return nil
}
