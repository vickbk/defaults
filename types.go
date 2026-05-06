package defaults

// Provider is a container for a fallback value of a specific type.
type Provider[T any] struct {
	defaultValue T
}

// Result encapsulates the outcome of a type check.
// It implements the error interface.
type Result struct {
	Message     string
	Ok          bool
	UsedDefault bool
}
