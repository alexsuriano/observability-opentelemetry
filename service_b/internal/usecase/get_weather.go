package usecase

import (
	"context"
	"net/http"
	"os"
	"strconv"

	externalapi "github.com/alexsuriano/observability-opentelemetry/service_b/internal/external-api"
	"go.opentelemetry.io/otel"
)

type GetWeatherHandler struct {
	Cep        string
	WeatherApi externalapi.WeatherApiHandler
	ViaCep     externalapi.ViaCepHandler
}

type Temperature struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type ErrorResponse struct {
	HttpCode int    `json:"HttpCode"`
	Message  string `json:"message"`
}

func NewGetWeatherHandler(
	cep string,
	weatherApi externalapi.WeatherApiHandler,
	viaCep externalapi.ViaCepHandler,
) *GetWeatherHandler {
	return &GetWeatherHandler{
		Cep:        cep,
		WeatherApi: weatherApi,
		ViaCep:     viaCep,
	}
}

func (g *GetWeatherHandler) Execute(ctx context.Context) (*Temperature, *ErrorResponse) {
	if !g.ValidateCep() {
		return nil, &ErrorResponse{
			HttpCode: http.StatusUnprocessableEntity,
			Message:  "invalid zipcode",
		}
	}

	ctx, spanViaCep := otel.Tracer(os.Getenv("OTEL_SERVICE_NAME")).Start(ctx, "Via-CEP")
	location, err := g.ViaCep.GetLocation(g.Cep)
	spanViaCep.End()
	if err != nil {
		if err.Error() == "can not find zipcode" {
			return nil, &ErrorResponse{
				HttpCode: http.StatusNotFound,
				Message:  "can not find zipcode",
			}
		} else {
			return nil, &ErrorResponse{
				HttpCode: http.StatusInternalServerError,
				Message:  err.Error(),
			}
		}
	}
	ctx, spanWeatherApi := otel.Tracer(os.Getenv("OTEL_SERVICE_NAME")).Start(ctx, "Weather-API")
	weather, err := g.WeatherApi.GetTemperature(location.Localidade)
	spanWeatherApi.End()
	if err != nil {
		return nil, &ErrorResponse{
			HttpCode: http.StatusBadRequest,
			Message:  err.Error(),
		}
	}

	return &Temperature{
		City:  weather.Location.Name,
		TempC: weather.Current.TempC,
		TempF: (weather.Current.TempC * 1.8) + 32,
		TempK: weather.Current.TempC + 273.15,
	}, nil
}

func (g *GetWeatherHandler) ValidateCep() bool {
	if g.Cep != "" {
		if len(g.Cep) == 8 {
			_, err := strconv.Atoi(g.Cep)
			return err == nil
		}
	}

	return false
}
