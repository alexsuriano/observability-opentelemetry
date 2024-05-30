package web

import (
	"encoding/json"
	"net/http"
	"os"
	"regexp"

	externalapi "github.com/alexsuriano/observability-opentelemetry/service_a/internal/external-api"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type InputGetTemp struct {
	CEP string `json:"cep"`
}

type ResponseGetTemp struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
}

type GetTempHandler struct {
	Tracer trace.Tracer
}

func NewGetTempHandler(tracer trace.Tracer) *GetTempHandler {
	return &GetTempHandler{
		Tracer: tracer,
	}
}

func (g *GetTempHandler) GetTemp2(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	carrier := propagation.HeaderCarrier(r.Header)
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := g.Tracer.Start(ctx, "service-a-get-temp")
	defer span.End()

	var input InputGetTemp

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if validateCep(input.CEP) == false {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	serviceB := externalapi.NewServiceB()
	serviceBResponse, errResponse := serviceB.GetCityTemp(ctx, input.CEP)
	if errResponse != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, errResponse.Message, errResponse.HttpCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(serviceBResponse)

}

func GetTemp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	carrier := propagation.HeaderCarrier(r.Header)
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	ctx, span := otel.Tracer(os.Getenv("OTEL_SERVICE_NAME")).Start(ctx, "GetTemp")
	defer span.End()

	var input InputGetTemp

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if validateCep(input.CEP) == false {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	serviceB := externalapi.NewServiceB()
	serviceBResponse, errResponse := serviceB.GetCityTemp(ctx, input.CEP)
	if errResponse != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, errResponse.Message, errResponse.HttpCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(serviceBResponse)

}

func validateCep(cep string) bool {
	regex := `^\d{8}$`
	pattern := regexp.MustCompile(regex)
	return pattern.MatchString(cep)
}
