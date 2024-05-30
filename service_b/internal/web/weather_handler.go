package web

import (
	"encoding/json"
	"net/http"

	externalapi "github.com/alexsuriano/observability-opentelemetry/service_b/internal/external-api"
	"github.com/alexsuriano/observability-opentelemetry/service_b/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func GetTemp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	carrier := propagation.HeaderCarrier(r.Header)
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	ctx, span := otel.Tracer("service_b").Start(ctx, "weather_handler")
	defer span.End()

	cepParam := r.URL.Query().Get("cep")
	viaCep := externalapi.NewViaCepHandler()
	weatherApi := externalapi.NewWeatherApiHandler()

	getWeather := usecase.NewGetWeatherHandler(
		cepParam,
		*weatherApi,
		*viaCep,
	)

	temp, err := getWeather.Execute(ctx)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, err.Message, err.HttpCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(temp)

}
