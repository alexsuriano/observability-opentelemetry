package webserver

import (
	"log"
	"net/http"
)

type WebServer struct {
	WebServerPort string
	Mux           *http.ServeMux
	Handlers      map[string]http.HandlerFunc
}

func NewWebServer(port string, mux *http.ServeMux) *WebServer {
	return &WebServer{
		WebServerPort: port,
		Mux:           mux,
		Handlers:      make(map[string]http.HandlerFunc),
	}
}

func (s *WebServer) AddHandler(pattern string, handler http.HandlerFunc) {
	s.Mux.HandleFunc(pattern, handler)
}

func (s *WebServer) Start() {
	server := http.Server{
		Addr:    s.WebServerPort,
		Handler: s.Mux,
	}

	log.Printf("Starting web server on port: %s", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
