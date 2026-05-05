package defaults

import (
	"fmt"
	"reflect"
)

// Check validates if an input matches type T.
// It returns the defaultValue and UsedDefault=true if the input is nil or a typed nil pointer.
// If the type is mismatched, it returns an error status with Ok=false.
// The optional message parameter allows for custom error messages on type mismatch.
// Note: Be careful when using this directly with variadic arguments, as it may lead to "index out of range" panics if the expected index is not provided.
//
//	Consider using SafeCheck for safer access to variadic arguments.
//	Or can use DefaultArgsNormalize to ensure the slice has enough elements before calling Check.
//
// Example:
//
//	val, status := DefaultValue(10).Check("not an int", "Age must be a number")
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

// SafeCheck performs an index-safe check on a slice.
// If the index is out of bounds, it automatically returns the default value.
//
// Example:
//
//	val, status := DefaultValue("Guest").SafeCheck(vars, 0)
func (d Default[T]) SafeCheck(values []any, index int, message ...string) (T, DefaultType) {
	value := any(nil)
	if index < len(values) {
		value = values[index]
	}

	return d.Check(value, message...)
}

// SafeCheckOrPanic behaves like SafeCheck but panics if a type mismatch occurs.
// Use this for critical internal configurations that must be of a specific type.
//
// Example:
//
//	port := DefaultValue(8080).SafeCheckOrPanic(args, 0)
func (d Default[T]) SafeCheckOrPanic(values []any, index int, message ...string) T {
	val, status := d.SafeCheck(values, index, message...)
	if !status.Ok {
		panic(status.Message)
	}
	return val
}
