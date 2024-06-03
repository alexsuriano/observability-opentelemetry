package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alexsuriano/observability-opentelemetry/service_a/internal/web"
	"github.com/alexsuriano/observability-opentelemetry/service_a/internal/web/webserver"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	resource2 "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func initProvider() (func(context.Context) error, error) {
	ctx := context.Background()

	collectorUL := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")

	res, err := resource2.New(ctx,
		resource2.WithAttributes(
			semconv.ServiceName(os.Getenv("OTEL_SERVICE_NAME"))))
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, collectorUL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(bsp),
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
