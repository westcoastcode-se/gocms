package jwt

import "github.com/westcoastcode-se/gocms/pkg/security"

// Interface responsible for converting a user back and forth between a token and an object
type Tokenizer interface {
	// Generate a token based on the supplied user
	UserToToken(user *security.User) (string, error)

	// Create a user object based on the supplied token
	TokenToUser(token string) (*security.User, error)
}
