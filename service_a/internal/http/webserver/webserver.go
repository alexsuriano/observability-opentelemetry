package webserver

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	WebServerPor string
	Router       chi.Router
	Handlers     map[string]http.HandlerFunc
}

func NewWebServer(port string) *WebServer {
	return &WebServer{
		WebServerPor: port,
		Router:       chi.NewRouter(),
		Handlers:     make(map[string]http.HandlerFunc),
	}
}

func (s *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)

	for path, handler := range s.Handlers {
		s.Router.Handle(path, handler)
	}

	server := http.Server{
		Addr:    s.WebServerPor,
		Handler: s.Router,
	}

	log.Printf("Starting web server on port: %s", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
