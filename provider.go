package defaults

// Check validates if an input matches type T.
// It returns the defaultValue and UsedDefault=true if the input is nil or a typed nil pointer.
// If the type is mismatched, it returns an error status with Ok=false.
// The optional message parameter allows for custom error messages on type mismatch.
// Note: Be careful when using this directly with variadic arguments, as it may lead to "index out of range" panics if the expected index is not provided (Normal case of optional arguments).
//
//	Consider using SafeCheck for safer access to variadic arguments.
//	Or can use Normalize to ensure the slice has enough elements before calling Check.
//
// Example:
//
//	val, status := defaults.Value(10).Check("not an int", "Age must be a number")
//
// Deprecated: use defaults.Safe instead but make sure to pass in the entire slice insted of indexing it directly
func (d Provider[T]) Check(input any, message ...string) (T, Result) {
	return Safe([]any{input}, d.defaultValue, message...)
}

// SafeCheck performs an index-safe check on a slice.
// If the index is out of bounds, it automatically returns the default value.
//
// Example:
//
//	val, status := defaults.Value("Guest").SafeCheck(vars, 0)
//
// Deprecated: use defaults.SafeAt instead for improved performance
func (d Provider[T]) SafeCheck(values []any, index int, message ...string) (T, Result) {
	return SafeAt(values, index, d.defaultValue, message...)
}

// SafeCheckOrPanic behaves like SafeCheck but panics if a type mismatch occurs.
// Use this for critical internal configurations that must be of a specific type.
//
// Example:
//
//	port := defaults.Value(8080).SafeCheckOrPanic(args, 0)
//
// Deprecated: use defaults.Required instead for better performance
func (d Provider[T]) SafeCheckOrPanic(values []any, index int, message ...string) T {
	return Required(values, index, d.defaultValue, message...)
}
