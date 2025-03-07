package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

// LoggerKey is the context key for the request logger
type LoggerKey struct{}

// LogHandler returns a middleware that adds a logger to the request context
// and logs each request with its details (method, URL, status, duration)
func LogHandler(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add logger to request context
			ctx := r.Context()
			reqLogger := logger.With("req_id", r.Header.Get("X-Request-ID"))
			ctx = context.WithValue(ctx, LoggerKey{}, reqLogger)
			r = r.WithContext(ctx)

			// Create response writer wrapper to capture status code
			rw := responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Record timing
			startTime := time.Now()
			next.ServeHTTP(&rw, r)
			duration := time.Since(startTime)

			// Log request details
			reqLogger.LogAttrs(ctx, slog.LevelInfo,
				"HTTP Request",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.Int("status", rw.statusCode),
				slog.Duration("duration", duration),
			)
		})
	}
}

// LoggerFromRequest extracts the logger from the request context
func LoggerFromRequest(r *http.Request) *slog.Logger {
	logger, ok := r.Context().Value(LoggerKey{}).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return logger
}

// responseWriter is a wrapper around http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write captures a 200 status code if WriteHeader hasn't been called
func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	return rw.ResponseWriter.Write(b)
}
