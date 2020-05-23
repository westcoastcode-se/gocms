package content

// Error raised when a model at a specific path is not found. This normally results in a HTTP 404.
type NotFoundError struct {
	message string
}

func (p *NotFoundError) Error() string {
	return p.message
}

// Create a new NotFoundError
func NewNotFoundError(path string) *NotFoundError {
	return &NotFoundError{message: "could not find: " + path}
}
