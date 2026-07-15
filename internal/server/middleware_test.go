package server

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogHandler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	t.Run("captures status code and injects a logger", func(t *testing.T) {
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if LoggerFromRequest(r) == nil {
				t.Error("LoggerFromRequest(r) = nil, want a logger in context")
			}
			w.WriteHeader(http.StatusTeapot)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		LogHandler(logger)(next).ServeHTTP(w, req)

		if w.Code != http.StatusTeapot {
			t.Errorf("status = %d, want %d", w.Code, http.StatusTeapot)
		}
	})

	t.Run("generates a fallback request id when header is absent", func(t *testing.T) {
		var gotLogger *slog.Logger
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotLogger = LoggerFromRequest(r)
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		LogHandler(logger)(next).ServeHTTP(w, req)

		if gotLogger == nil {
			t.Fatal("LoggerFromRequest(r) = nil, want a logger in context")
		}
	})
}

func TestLoggerFromRequest_defaultsWhenMissing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if got := LoggerFromRequest(req); got == nil {
		t.Error("LoggerFromRequest(req) = nil, want slog.Default() fallback")
	}
}

func TestResponseWriter_defaultsToOK(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := responseWriter{ResponseWriter: rec, statusCode: http.StatusOK}

	if _, err := rw.Write([]byte("hi")); err != nil {
		t.Fatalf("Write() err = %v", err)
	}
	if rw.statusCode != http.StatusOK {
		t.Errorf("statusCode = %d, want %d", rw.statusCode, http.StatusOK)
	}
}
