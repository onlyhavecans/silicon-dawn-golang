package server

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/onlyhavecans/silicondawn/internal/cards"

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

// NewServer returns an initialized Server
func NewServer(config *Config) *Server {
	deck, err := cards.NewCardDeck(config.CardsDir)
	if err != nil {
		slog.Error("deck build failure", slog.Any("error", err))
		os.Exit(1)
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

	templates := template.Must(template.ParseGlob(templatesPath))

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

	return server
}

// Start starts the httpServer & returns when done
func (s *Server) Start() error {
	s.logger.LogAttrs(context.Background(), slog.LevelInfo, "Listening",
		slog.String("port", s.config.Port),
		slog.Int("card count", s.deck.Count()),
	)
	return s.httpServer.ListenAndServe()
}

func (s *Server) robots(w http.ResponseWriter, _ *http.Request) {
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

	w.Header().Set("Content-type", "text/html")

	err = s.templates.ExecuteTemplate(w, templateIndex, map[string]string{
		"dir":  "cards",
		"name": c.Front(),
		"text": c.Back(),
	})
	if err != nil {
		LoggerFromRequest(r).LogAttrs(r.Context(), slog.LevelWarn, "template render error",
			slog.Any("error", err))
		s.errorHandler(w, r, http.StatusInternalServerError)
	}
}

func (s *Server) errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-type", "text/html")

	err := s.templates.ExecuteTemplate(w, templateError, map[string]string{
		"status": strconv.Itoa(status),
	})
	if err != nil {
		LoggerFromRequest(r).LogAttrs(r.Context(), slog.LevelError, "could not render error page",
			slog.Any("error", err))
	}
}
