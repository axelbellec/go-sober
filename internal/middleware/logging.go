package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type LoggingMiddleware struct{}

func NewLoggingMiddleware() *LoggingMiddleware {
	return &LoggingMiddleware{}
}

func (m *LoggingMiddleware) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a wrapper for the response writer
		wrapper := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default to 200 OK
		}

		next.ServeHTTP(wrapper, r)

		slog.Info("Request processed",
			"method", r.Method,
			"path", r.RequestURI,
			"duration", time.Since(start),
			"status", wrapper.statusCode,
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
