package usecase

import (
	"testing"

	externalapi "github.com/alexsuriano/observability-opentelemetry/service_b/internal/external-api"
	"github.com/stretchr/testify/assert"
)

func TestGetWeatherExecute(t *testing.T) {
	var cep = "17018550"
	weatherApi := externalapi.NewWeatherApiHandler()
	viaCepApi := externalapi.NewViaCepHandler()

	getWeather := NewGetWeatherHandler(cep, *weatherApi, *viaCepApi)

	temp, err := getWeather.Execute()

	assert.Nil(t, err)
	assert.NotNil(t, temp)
}

func TestValidateCep(t *testing.T) {

	t.Run("When zipcode is valid, should return true", func(t *testing.T) {
		var cep = "01000000"
		weatherApi := externalapi.NewWeatherApiHandler()
		viaCepApi := externalapi.NewViaCepHandler()

		getWeather := NewGetWeatherHandler(cep, *weatherApi, *viaCepApi)

		valid := getWeather.ValidateCep()
		assert.True(t, valid)
	})

	t.Run("When zipcode is short, should return false", func(t *testing.T) {
		var cep = "1234567"
		weatherApi := externalapi.NewWeatherApiHandler()
		viaCepApi := externalapi.NewViaCepHandler()

		getWeather := NewGetWeatherHandler(cep, *weatherApi, *viaCepApi)

		valid := getWeather.ValidateCep()
		assert.False(t, valid)

	})

	t.Run("When zipcode is invalid, should return false", func(t *testing.T) {
		var cep = "0100abc0"
		weatherApi := externalapi.NewWeatherApiHandler()
		viaCepApi := externalapi.NewViaCepHandler()

		getWeather := NewGetWeatherHandler(cep, *weatherApi, *viaCepApi)

		valid := getWeather.ValidateCep()
		assert.False(t, valid)

	})
}
