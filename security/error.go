package security

// Error raised when a user is not found
type UserNotFound struct {
	Username string
}

func (u *UserNotFound) Error() string {
	return `user "` + u.Username + `" is not found"`
}
