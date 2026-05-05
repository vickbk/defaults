package defaults

// Error allows DefaultType to satisfy the error interface for use in CheckDefaults.
func (d DefaultType) Error() string {
	return d.Message
}
