package acl

import "fmt"

// Error raised when an error occurs during a load of some kind. For example when a database is loaded
type LoadError struct {
	message string
}

func (l *LoadError) Error() string {
	return l.message
}

func NewLoadError(format string, v ...interface{}) *LoadError {
	return &LoadError{
		message: fmt.Sprintf(format, v),
	}
}
