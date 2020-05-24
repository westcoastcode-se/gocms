package security

import "fmt"

// Error raised when a user is not found
type UserNotFound struct {
	Username string
}

func (u *UserNotFound) Error() string {
	return `user "` + u.Username + `" is not found"`
}

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
