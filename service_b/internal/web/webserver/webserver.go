package webserver

import (
	"log"
	"net/http"
)

type Webserver struct {
	WebserverPort string
	Mux           *http.ServeMux
	Handlers      map[string]http.HandlerFunc
}

func NewWebserver(port string, mux *http.ServeMux) *Webserver {
	return &Webserver{
		WebserverPort: port,
		Mux:           mux,
		Handlers:      make(map[string]http.HandlerFunc),
	}
}

func (s *Webserver) AddHandler(pattern string, handler http.HandlerFunc) {
	s.Mux.HandleFunc(pattern, handler)
}

func (s *Webserver) Start() {
	server := http.Server{
		Addr:    s.WebserverPort,
		Handler: s.Mux,
	}

	log.Printf("Starting web server on port: %s", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
