package auth

import "github.com/westcoastcode-se/gocms/pkg/security"

// Service used when looking up users
type LoginService interface {
	// Try to login with the supplied username and password
	Login(username string, password string) (*security.User, error)
}
