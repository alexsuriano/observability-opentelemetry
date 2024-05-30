package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alexsuriano/observability-opentelemetry/service_a/internal/web"
	"github.com/alexsuriano/observability-opentelemetry/service_a/internal/web/webserver"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	resource2 "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func initProvider() (func(context.Context) error, error) {
	ctx := context.Background()

	exporter, err := zipkin.New(os.Getenv("ZIPKIN_ENDPOINT"))
	if err != nil {
		log.Fatal("Fail to create zipkin exporter: %v", err)
	}

	res, err := resource2.New(ctx, resource2.WithAttributes(
		semconv.ServiceName(os.Getenv("OTEL_SERVICE_NAME"))))

	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider.Shutdown, nil
}

func main() {
	initProvider()

	mux := http.NewServeMux()
	server := webserver.NewWebServer(":8181", mux)

	server.AddHandler("POST /temperature", web.GetTemp)
	server.AddHandler("GET /metrics", promhttp.Handler().ServeHTTP)

	server.Start()
}
