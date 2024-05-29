package externalapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTemperature(t *testing.T) {
	var location = "Bauru"
	weatherApi := NewWeatherApiHandler()

	response, err := weatherApi.GetTemperature(location)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Current.TempC)
}
