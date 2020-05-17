package security

import "github.com/dgrijalva/jwt-go"

const (
	Read       = "Read"
	Write      = "Write"
	Admin      = "Admin"
	SessionKey = "X-Auth-Token"
	Issuer     = "gocms"
)

type Claims struct {
	Name  string
	Roles []string
	jwt.StandardClaims
}
