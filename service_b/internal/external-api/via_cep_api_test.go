package externalapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLocation(t *testing.T) {
	var cep = "17018550"
	var expectedLocation = "Bauru"
	viaCep := NewViaCepHandler()

	response, err := viaCep.GetLocation(cep)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedLocation, response.Localidade)
}
