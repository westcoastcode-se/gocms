package middleware

import (
	"context"
	"github.com/westcoastcode-se/gocms/pkg/security"
	"log"
	"net/http"
)

func Security(tokenizer security.Tokenizer, next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie(security.SessionKey)
		var user *security.User
		if err == nil && token.Value != "" {
			user, err = tokenizer.TokenToUser(token.Value)
			if err != nil {
				log.Printf("User token could not be loaded. Reason %e\n", err)
			}
		}
		if user == nil {
			user = security.NotLoggedInUser
		}
		r = r.WithContext(context.WithValue(r.Context(), security.SessionKey, user))
		next.ServeHTTP(rw, r)
	})
}
