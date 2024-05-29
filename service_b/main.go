package main

import (
	"github.com/alexsuriano/observability-opentelemetry/service_b/internal/http"
	"github.com/alexsuriano/observability-opentelemetry/service_b/internal/http/webserver"
)

func main() {
	server := webserver.NewWebserver(":8282")
	http.SetRoutes(server)
	server.Start()
}
