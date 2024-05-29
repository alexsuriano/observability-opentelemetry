package externalapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type ViaCepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type ViaCepHandler struct{}

func NewViaCepHandler() *ViaCepHandler {
	return &ViaCepHandler{}
}

func (v *ViaCepHandler) GetLocation(cep string) (*ViaCepResponse, error) {
	var viaCepResponse ViaCepResponse
	url := generateUrlViaCep(cep)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var resp map[string]interface{}
		err := json.Unmarshal(body, &resp)
		if err != nil {
			return nil, err
		}

		_, ok := resp["erro"]
		if ok {
			return nil, errors.New("can not find zipcode")
		}

		err = json.Unmarshal(body, &viaCepResponse)
		if err != nil {
			return nil, err
		}

		return &viaCepResponse, nil
	}

	return nil, errors.New("unexpected response")
}

func generateUrlViaCep(cep string) string {
	return fmt.Sprintf("http://viacep.com.br/ws/%s/json", cep)
}
