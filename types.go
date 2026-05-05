package defaults

// Default holds the fallback value for a specific type.
type Default[T any] struct {
	defaultValue T
}

type DefaultType struct {
	Message string
	Ok      bool
}
