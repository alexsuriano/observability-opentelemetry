package main

import (
	"github.com/alexsuriano/observability-opentelemetry/service_a/internal/http"
	"github.com/alexsuriano/observability-opentelemetry/service_a/internal/http/webserver"
)

func main() {

	server := webserver.NewWebServer(":8181")
	http.SetRoutes(server)
	server.Start()
}
