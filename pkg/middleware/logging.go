package middleware

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func getIpAddress(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		logger := logrus.WithField("id", uuid.New().String())
		start := time.Now()
		defer func() {
			diff := time.Since(start)
			logger.WithFields(logrus.Fields{
				"uri":     r.RequestURI,
				"method":  r.Method,
				"remote":  getIpAddress(r),
				"elapsed": diff,
			}).Info()
		}()
		next.ServeHTTP(rw, r)
	})
}
