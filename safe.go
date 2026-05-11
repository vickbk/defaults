package defaults

// defaults.Safe validates if the zero index value type matches type T of default value.
// It returns the defaultValue and UsedDefault=true if the input is nil or a typed nil pointer.
// If the type is mismatched, it returns an error status with Ok=false.
// The optional message parameter allows for custom error messages on type mismatch.
// Note: Be careful when using Safe only for one optional parameter as in introduces type overhead.
// Use instead defaults.Get for only one optional parameter for type safety.
//
// Example:
//
//	val, status := defaults.Safe(args, 10, "Age must be an int")
func Safe[T any](values []any, defaultValue T, message ...string) (T, Result) {
	return SafeAt(values, 0, defaultValue, message...)
}
