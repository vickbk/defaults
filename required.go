package defaults

// Required behaves like SafeAt but panics if a type mismatch occurs.
// Use this for critical internal configurations that must be of a specific type.
// Panics on type mismatch with the provided error message.
//
// Example:
//
//	port := Required(args, 0, 8080, "port must be int")
//	// Panics if args[0] is not an int or is a typed nil
func Required[T any](values []any, index int, defaultValue T, message ...string) T {
	val, status := SafeAt(values, index, defaultValue, message...)
	if !status.Ok {
		panic(status.Message)
	}
	return val
}
