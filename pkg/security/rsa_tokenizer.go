package security

import (
	"crypto/rsa"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"time"
)

type RsaTokenizer struct {
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
}

func (r *RsaTokenizer) UserToToken(user *User) (string, error) {
	now := time.Now()
	expirationTime := now.Add(5 * time.Minute)
	claims := &Claims{
		Name:  user.Name,
		Roles: user.Roles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    Issuer,
			IssuedAt:  now.Unix(),
			Subject:   user.Name,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(r.signKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (r *RsaTokenizer) TokenToUser(tokenAsString string) (*User, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenAsString, claims, func(token *jwt.Token) (interface{}, error) {
		return r.verifyKey, nil
	})

	if err != nil {
		return NotLoggedInUser, errors.New("could not parse token claims. Reason: " + err.Error())
	}

	if !token.Valid {
		return NotLoggedInUser, errors.New("token is no longer valid")
	}

	return &User{claims.Subject, claims.Roles}, nil
}

// Create a new asymmetric tokenizer instance used.
func NewAsymmetricTokenizer(publicKey string, privateKey string) Tokenizer {
	verifyBytes, err := ioutil.ReadFile(publicKey)
	if err != nil {
		panic(err)
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		panic(err)
	}

	signBytes, err := ioutil.ReadFile(privateKey)
	if err != nil {
		panic(err)
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		panic(err)
	}

	return &RsaTokenizer{
		verifyKey: verifyKey,
		signKey:   signKey,
	}
}
