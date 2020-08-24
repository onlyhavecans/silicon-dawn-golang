package silicondawn

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

//Config stores the server's config
type Config struct {
	Port     string
	CardsDir string
}

//Server is an instance of silicon-dawn server
type Server struct {
	httpServer *http.Server
	config     *Config
	templates  *template.Template
	deck       *CardDeck
}

const rootPath = "/"
const cardsPath = "/cards/"

const indexTemplatePath = "templates/index.gohtml"

//NewServer returns an initialized Server
func NewServer(config *Config) *Server {
	deck, err := NewCardDeck(config.CardsDir)
	if err != nil {
		log.Fatalf("deck build error: %v", err)
	}

	mux := http.NewServeMux()

	httpServer := &http.Server{
		Addr:         ":" + config.Port,
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
		IdleTimeout:  1 * time.Minute,
		Handler:      mux,
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

func (s *Server) Start() error {
	log.Printf("Listening on %s with %d cards", s.config.Port, s.deck.Count())
	return s.httpServer.ListenAndServe()
}

func (s *Server) robots(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("User-agent: *\nDisallow: /\n"))
}

func (s *Server) root(w http.ResponseWriter, r *http.Request) {
	// I prefer 404
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	c, err := s.deck.Draw()
	if err != nil {
		err := fmt.Sprintf("could not draw card: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err))
		return
	}

	w.Header().Set("Content-type", "text/html")
	err = s.templates.ExecuteTemplate(w, "index.gohtml", map[string]string{
		"dir":  "cards",
		"name": c.Front(),
		"text": c.Back(),
	})
	if err != nil {
		log.Printf("Error executing template rendering: %v", err)
	}
}
