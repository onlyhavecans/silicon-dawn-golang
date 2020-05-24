package silicondawn

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"onlyhavecans.works/amy/silicondawn/lib"
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
}

const rootPath = "/"
const cardsPath = "/cards/"

const indexTemplatePath = "templates/index.gohtml"

var deck lib.CardDeck

func initDeck(cardsDirectory string) error {
	log.Printf("Building a deck out of %s", cardsDirectory)
	var err error
	deck, err = lib.NewCardDeck(cardsDirectory)
	if err != nil {
		return fmt.Errorf("building deck: %w", err)
	}
	log.Printf("We have %d cards now", deck.Count())
	return nil
}

//NewServer returns an initialized Server
func NewServer(config *Config) *Server {
	err := initDeck(config.CardsDir)
	if err != nil {
		log.Fatal(err)
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
	}

	mux.HandleFunc("/robots.txt", server.robots)
	mux.HandleFunc(rootPath, server.root)
	c := http.Dir(config.CardsDir)
	mux.Handle(cardsPath, http.StripPrefix(cardsPath, http.FileServer(c)))

	return server
}

func (s *Server) Start() error {
	log.Printf("Listening on %v", s.config.Port)

	return s.httpServer.ListenAndServe()
}

func (s Server) robots(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("User-agent: *\nDisallow: /\n"))
}

func (s Server) root(w http.ResponseWriter, r *http.Request) {
	// I prefer 404
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	c, err := deck.Draw()
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
