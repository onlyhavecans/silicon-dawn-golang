package server

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"onlyhavecans.works/onlyhavecans/silicon-dawn/internal/cards"

	"github.com/justinas/alice"
)

// Config stores the server's config
type Config struct {
	Port     string
	CardsDir string
	LogLevel slog.Level
}

// Server is an instance of silicon-dawn server
type Server struct {
	httpServer *http.Server
	config     *Config
	templates  *template.Template
	deck       *cards.CardDeck
	logger     *slog.Logger
}

const (
	rootPath  = "/"
	cardsPath = "/cards/"
)

const (
	templatesPath = "templates/*"
	templateIndex = "index.gohtml"
	templateError = "error.gohtml"
)

// NewServer returns an initialized Server, or an error if the card deck or
// templates could not be loaded.
func NewServer(config *Config) (*Server, error) {
	deck, err := cards.NewCardDeck(config.CardsDir)
	if err != nil {
		return nil, fmt.Errorf("deck build failure: %w", err)
	}

	templates, err := template.ParseGlob(templatesPath)
	if err != nil {
		return nil, fmt.Errorf("template parse failure: %w", err)
	}

	// Configure logger
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: config.LogLevel,
	})
	logger := slog.New(handler)

	mux := http.NewServeMux()

	// Setup middleware chain
	a := alice.New()
	a = a.Append(LogHandler(logger))
	h := a.Then(mux)

	httpServer := &http.Server{
		Addr:         ":" + config.Port,
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
		IdleTimeout:  1 * time.Minute,
		Handler:      h,
	}

	server := &Server{
		httpServer: httpServer,
		config:     config,
		templates:  templates,
		deck:       deck,
		logger:     logger,
	}

	mux.HandleFunc("/robots.txt", server.robots)
	mux.HandleFunc(rootPath, server.root)
	c := http.Dir(config.CardsDir)
	mux.Handle(cardsPath, http.StripPrefix(cardsPath, http.FileServer(c)))

	return server, nil
}

// Start starts the httpServer and blocks until ctx is canceled, at which
// point it gracefully shuts down in-flight requests before returning.
func (s *Server) Start(ctx context.Context) error {
	s.logger.LogAttrs(ctx, slog.LevelInfo, "Listening",
		slog.String("port", s.config.Port),
		slog.Int("card count", s.deck.Count()),
	)

	serveErr := make(chan error, 1)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serveErr <- err
			return
		}
		serveErr <- nil
	}()

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
		s.logger.LogAttrs(context.Background(), slog.LevelInfo, "shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}
		return nil
	}
}

func (s *Server) robots(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-type", "text/plain")
	_, _ = w.Write([]byte("User-agent: *\nDisallow: /\n"))
}

func (s *Server) root(w http.ResponseWriter, r *http.Request) {
	// I prefer 404
	if r.URL.Path != "/" {
		s.errorHandler(w, r, http.StatusNotFound)
		return
	}

	c, err := s.deck.Draw()
	if err != nil {
		LoggerFromRequest(r).LogAttrs(r.Context(), slog.LevelError, "could not draw card",
			slog.Any("error", err))
		s.errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = s.templates.ExecuteTemplate(&buf, templateIndex, map[string]string{
		"dir":  "cards",
		"name": c.Front(),
		"text": c.Back(),
	})
	if err != nil {
		LoggerFromRequest(r).LogAttrs(r.Context(), slog.LevelWarn, "template render error",
			slog.Any("error", err))
		s.errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "text/html")
	_, _ = io.Copy(w, &buf)
}

func (s *Server) errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	var buf bytes.Buffer
	err := s.templates.ExecuteTemplate(&buf, templateError, map[string]string{
		"status": strconv.Itoa(status),
	})
	if err != nil {
		LoggerFromRequest(r).LogAttrs(r.Context(), slog.LevelError, "could not render error page",
			slog.Any("error", err))
		w.WriteHeader(status)
		return
	}

	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(status)
	_, _ = io.Copy(w, &buf)
}
