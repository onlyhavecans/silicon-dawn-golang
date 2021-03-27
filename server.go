package silicondawn

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

//Config stores the server's config
type Config struct {
	Port     string
	CardsDir string
	logLevel zerolog.Level
}

//Server is an instance of silicon-dawn server
type Server struct {
	httpServer *http.Server
	config     *Config
	templates  *template.Template
	deck       *CardDeck
	logger     zerolog.Logger
}

const rootPath = "/"
const cardsPath = "/cards/"

const indexTemplatePath = "templates/index.gohtml"

//NewServer returns an initialized Server
func NewServer(config *Config) *Server {
	deck, err := NewCardDeck(config.CardsDir)
	if err != nil {
		log.Fatal().Err(err).Msg("deck build failure")
	}

	mux := http.NewServeMux()

	l := zerolog.New(os.Stdout)
	l.Level(config.logLevel)

	a := alice.New()
	a = a.Append(hlog.NewHandler(l)).
		Append(hlog.AccessHandler(func(r *http.Request, status, _ int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Stringer("url", r.URL).
				Int("status", status).
				Dur("duration", duration).
				Msg("")
		}))
	h := a.Then(mux)

	httpServer := &http.Server{
		Addr:         ":" + config.Port,
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
		IdleTimeout:  1 * time.Minute,
		Handler:      h,
	}

	templates := template.Must(template.ParseFiles(
		indexTemplatePath,
	))

	server := &Server{
		httpServer: httpServer,
		config:     config,
		templates:  templates,
		deck:       deck,
	}

	mux.HandleFunc("/robots.txt", server.robots)
	mux.HandleFunc(rootPath, server.root)
	c := http.Dir(config.CardsDir)
	mux.Handle(cardsPath, http.StripPrefix(cardsPath, http.FileServer(c)))

	return server
}

// Start starts the httpServer & returns when done
func (s *Server) Start() error {
	log.Info().
		Str("port", s.config.Port).
		Int("card count", s.deck.Count()).
		Msg("Listening")
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
		log.Error().Err(err).Msg("could not draw card")
		s.errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "text/html")
	err = s.templates.ExecuteTemplate(w, "index.gohtml", map[string]string{
		"dir":  "cards",
		"name": c.Front(),
		"text": c.Back(),
	})
	if err != nil {
		log.Warn().Err(err).Msg("template render error")
		s.errorHandler(w, r, http.StatusInternalServerError)
	}
}

func (s *Server) errorHandler(w http.ResponseWriter, _ *http.Request, status int) {
	w.WriteHeader(status)
	msg := fmt.Sprintf("You drew a %d\nI doubt this was the card you are lookging for.", status)
	_, _ = w.Write([]byte(msg))
}
