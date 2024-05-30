package externalapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type ServiceBResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type ErrorResponse struct {
	HttpCode int    `json:"HttpCode"`
	Message  string `json:"message"`
}

type ServiceB struct{}

func NewServiceB() *ServiceB {
	return &ServiceB{}
}

func (s *ServiceB) GetCityTemp(ctx context.Context, cep string) (*ServiceBResponse, *ErrorResponse) {

	url := fmt.Sprintf("http://service_b:8282/temp?cep=%s", cep)

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, &ErrorResponse{
			HttpCode: http.StatusInternalServerError,
			Message:  err.Error(),
		}
	}

	request.Header.Set("Content-Type", "application/json")

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(request.Header))

	ctx, span := otel.Tracer(os.Getenv("OTEL_SERVICE_NAME")).Start(ctx, "GetCityTemp")
	response, err := http.DefaultClient.Do(request)
	span.End()
	if err != nil {
		return nil, &ErrorResponse{
			HttpCode: response.StatusCode,
			Message:  err.Error(),
		}
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, &ErrorResponse{
			HttpCode: http.StatusInternalServerError,
			Message:  err.Error(),
		}
	}

	if response.StatusCode != 200 {
		return nil, &ErrorResponse{
			HttpCode: response.StatusCode,
			Message:  string(body),
		}
	}

	var serviceBResponse *ServiceBResponse

	err = json.Unmarshal(body, &serviceBResponse)
	if err != nil {
		return nil, &ErrorResponse{
			HttpCode: http.StatusInternalServerError,
			Message:  err.Error(),
		}
	}

	return serviceBResponse, nil
}
