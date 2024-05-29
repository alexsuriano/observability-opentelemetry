package webserver

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Webserver struct {
	WebserverPort string
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
}

func NewWebserver(webserverport string) *Webserver {
	return &Webserver{
		WebserverPort: webserverport,
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]http.HandlerFunc),
	}
}

func (s *Webserver) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *Webserver) Start() {
	s.Router.Use(middleware.Logger)

	for path, handler := range s.Handlers {
		s.Router.Handle(path, handler)
	}

	server := http.Server{
		Addr:    s.WebserverPort,
		Handler: s.Router,
	}

	log.Printf("Starting service_b on port: %s", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
