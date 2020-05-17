package middleware

import (
	"bytes"
	"github.com/westcoastcode-se/gocms/cache"
	"net/http"
	"strconv"
)

type cachedResponseWriter struct {
	header     http.Header
	buffer     bytes.Buffer
	statusCode int
}

func (c *cachedResponseWriter) Header() http.Header {
	return c.header
}

func (c *cachedResponseWriter) Write(i []byte) (int, error) {
	return c.buffer.Write(i)
}

func (c *cachedResponseWriter) WriteHeader(statusCode int) {
	c.statusCode = statusCode
}

func Cache(cache cache.Pages, next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		uri := r.URL.Path

		// Ignore if caching is not allowed
		if !cache.IsAllowed(uri) {
			next.ServeHTTP(rw, r)
			return
		}

		// Return the cached result if found
		if cachedResult, err := cache.Find(uri); err == nil {
			contentLength := len(cachedResult)
			rw.Header().Set("Content-Length", strconv.Itoa(contentLength))
			rw.WriteHeader(http.StatusOK)
			_, _ = rw.Write(cachedResult)
			return
		}

		wrapper := &cachedResponseWriter{
			header:     make(http.Header),
			buffer:     bytes.Buffer{},
			statusCode: 0,
		}
		next.ServeHTTP(wrapper, r)

		b := wrapper.buffer.Bytes()
		if wrapper.statusCode < 400 {
			cache.Set(uri, b)
		}

		for key, values := range wrapper.header {
			for _, value := range values {
				rw.Header().Set(key, value)
			}
		}
		
		rw.Header().Set("Content-Length", strconv.Itoa(len(b)))
		rw.WriteHeader(wrapper.statusCode)
		_, _ = rw.Write(b)
	})
}
