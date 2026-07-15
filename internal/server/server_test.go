package server

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// newTestServer builds a Server rooted at the repo root, using the cards
// fixtures under testdata/cards and the real production templates.
func newTestServer(t *testing.T) *Server {
	t.Helper()
	t.Chdir("../..")

	srv, err := NewServer(&Config{
		Port:     "0",
		CardsDir: "internal/server/testdata/cards",
		LogLevel: slog.LevelError,
	})
	if err != nil {
		t.Fatalf("NewServer() err = %v", err)
	}
	return srv
}

func TestNewServer(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := newTestServer(t)
		if srv == nil {
			t.Fatal("NewServer() returned nil server with nil error")
		}
	})

	t.Run("bad cards dir", func(t *testing.T) {
		t.Chdir("../..")
		srv, err := NewServer(&Config{
			Port:     "0",
			CardsDir: "internal/server/testdata/does-not-exist",
			LogLevel: slog.LevelError,
		})
		if err == nil {
			t.Fatal("NewServer() err = nil, want error for missing cards dir")
		}
		if srv != nil {
			t.Fatalf("NewServer() server = %v, want nil on error", srv)
		}
	})
}

func TestServer_root(t *testing.T) {
	srv := newTestServer(t)

	t.Run("draws a card", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		srv.httpServer.Handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
		}
		if ct := w.Header().Get("Content-type"); ct != "text/html" {
			t.Errorf("Content-type = %q, want %q", ct, "text/html")
		}
		if !strings.Contains(w.Body.String(), "the-lovers.jpg") {
			t.Errorf("body = %q, want it to contain %q", w.Body.String(), "the-lovers.jpg")
		}
	})

	t.Run("unknown path is a 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/nope", nil)
		w := httptest.NewRecorder()

		srv.httpServer.Handler.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
		}
		if ct := w.Header().Get("Content-type"); ct != "text/html" {
			t.Errorf("Content-type = %q, want %q", ct, "text/html")
		}
		if !strings.Contains(w.Body.String(), "404") {
			t.Errorf("body = %q, want it to contain %q", w.Body.String(), "404")
		}
	})
}

func TestServer_robots(t *testing.T) {
	srv := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/robots.txt", nil)
	w := httptest.NewRecorder()

	srv.httpServer.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if ct := w.Header().Get("Content-type"); ct != "text/plain" {
		t.Errorf("Content-type = %q, want %q", ct, "text/plain")
	}
	if !strings.Contains(w.Body.String(), "Disallow") {
		t.Errorf("body = %q, want it to contain %q", w.Body.String(), "Disallow")
	}
}

func TestServer_errorHandler(t *testing.T) {
	srv := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	srv.errorHandler(w, req, http.StatusInternalServerError)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
	if ct := w.Header().Get("Content-type"); ct != "text/html" {
		t.Errorf("Content-type = %q, want %q", ct, "text/html")
	}
	if !strings.Contains(w.Body.String(), "500") {
		t.Errorf("body = %q, want it to contain %q", w.Body.String(), "500")
	}
}
