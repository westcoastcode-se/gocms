package middleware

import (
	"github.com/westcoastcode-se/gocms/security"
	"net/http"
	"net/url"
)

func Authorize(acl security.ACL, next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		uri := r.URL.Path
		user := r.Context().Value(security.SessionKey).(*security.User)
		roles := acl.GetRoles(uri)
		if !user.HasRoles(roles) {
			http.Redirect(rw, r, "/login?redirect="+url.QueryEscape(uri), http.StatusTemporaryRedirect)
			return
		}
		next.ServeHTTP(rw, r)
	})
}
