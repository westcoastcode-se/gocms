package middleware

import (
	"net/http"
)

func DefaultURI(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			r.URL.Path = "/index"
		}
		next.ServeHTTP(rw, r)
	})
}
