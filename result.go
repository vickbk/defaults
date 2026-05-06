package defaults

// Error allows Result to satisfy the error interface.
func (d Result) Error() string {
	return d.Message
}
