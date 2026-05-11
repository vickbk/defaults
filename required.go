package defaults

// defaults.Required behaves like default.SafeAt but panics if a type mismatch occurs.
// Use this for critical internal configurations that must be of a specific type.
//
// Example:
//
//	port := defaults.Required(args, 0, 8080)
func Required[T any](values []any, index int, defaultValue T, message ...string) T {
	val, status := SafeAt(values, index, defaultValue, message...)
	if !status.Ok {
		panic(status.Message)
	}
	return val
}
