package cms

import "net/http"

func PreProcessURI(r *http.Request) string {
	if r.URL.Path == "/" {
		r.URL.Path = "/index"
	}
	return r.URL.Path
}
