package jwt

import "github.com/dgrijalva/jwt-go"

const (
	SessionKey = "X-Auth-Token"
	Issuer     = "gocms"
)

type Claims struct {
	Name  string
	Roles []string
	jwt.StandardClaims
}
