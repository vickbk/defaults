package defaults

// Default is a container for a fallback value of a specific type.
type Default[T any] struct {
	defaultValue T
}

// DefaultType encapsulates the result of a type check.
// It implements the error interface.
type DefaultType struct {
	Message     string
	Ok          bool
	UsedDefault bool
}
