package http

import "github.com/alexsuriano/observability-opentelemetry/service_b/internal/http/webserver"

func SetRoutes(server *webserver.Webserver) {
	server.AddHandler("/temp", GetTemp)
	server.AddHandler("/healthCheck", HealthCheck)
}
