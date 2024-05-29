package http

import (
	"encoding/json"
	"net/http"

	externalapi "github.com/alexsuriano/observability-opentelemetry/service_b/internal/external-api"
	"github.com/alexsuriano/observability-opentelemetry/service_b/internal/usecase"
)

func GetTemp(w http.ResponseWriter, r *http.Request) {
	cepParam := r.URL.Query().Get("cep")
	viaCep := externalapi.NewViaCepHandler()
	weatherApi := externalapi.NewWeatherApiHandler()

	getWeather := usecase.NewGetWeatherHandler(
		cepParam,
		*weatherApi,
		*viaCep,
	)

	temp, err := getWeather.Execute()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, err.Message, err.HttpCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(temp)

}
