package get_temperature

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aluferraz/go-expert-zipkin/internal/entity/zipcode"
	"github.com/aluferraz/go-expert-zipkin/internal/infra/http_clients"
	"io"
	"net/http"
	"net/url"
)

type AddressResponse struct {
	Localidade string `json:"localidade"`
}

type InputDTO struct {
	WeatherApiKey string
	WeatherApiUrl string
	ApiCepUrl     string
	Zipcode       zipcode.Zipcode
	Ctx           *context.Context
}
type OutputDTO struct {
	City           string  `json:"city"`
	CelsiusTemp    float64 `json:"temp_C"`
	FahrenheitTemp float64 `json:"temp_F"`
	KelvinTemp     float64 `json:"temp_K"`
}

type UseCase struct {
	client http_clients.ZipkinClientInterface
}

func NewUseCase(
	client http_clients.ZipkinClientInterface,
) UseCase {
	return UseCase{
		client: client,
	}
}

func (uc *UseCase) getCityFromZipCode(apiUrl string, zipcode string, ctx *context.Context) (string, error) {
	// "https://viacep.com.br/ws/%s/json/"

	newRequest, err := http.NewRequestWithContext(*ctx, "GET", fmt.Sprintf(apiUrl, zipcode), nil)
	if err != nil {
		return "", err
	}
	resp, err := uc.client.DoWithAppSpan(newRequest, "viacep")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var addressResponse AddressResponse
	err = json.Unmarshal(body, &addressResponse)
	if err != nil {
		return "", err
	}

	return addressResponse.Localidade, nil
}

func (uc *UseCase) getTemperatureFromCity(apiUrl string, apikey string, city string, ctx *context.Context) (float64, error) {
	//"http://api.weatherapi.com/v1/current.json?key=%s&q=%s"
	url := fmt.Sprintf(apiUrl, url.QueryEscape(apikey), url.QueryEscape(city))

	newRequest, err := http.NewRequestWithContext(*ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := uc.client.DoWithAppSpan(newRequest, "weather_api")
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		Body.Close()
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	var weatherResponse map[string]interface{}
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		return 0, err
	}

	temperatureData, ok := weatherResponse["current"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("Invalid temperature data")
	}

	return temperatureData["temp_c"].(float64), nil
}

func (uc *UseCase) Execute(inputDto InputDTO) (OutputDTO, error) {
	var err error
	city, err := uc.getCityFromZipCode(inputDto.ApiCepUrl, inputDto.Zipcode.Zipcode, inputDto.Ctx)
	result := OutputDTO{}
	if err != nil {
		return result, err
	}
	if city != "" {
		celsiusTemp, err := uc.getTemperatureFromCity(inputDto.WeatherApiUrl, inputDto.WeatherApiKey, city, inputDto.Ctx)
		if err != nil {
			return result, err
		}
		result.CelsiusTemp = celsiusTemp
	}

	result.City = city
	result.FahrenheitTemp = (result.CelsiusTemp * 1.8) + 32
	result.KelvinTemp = result.CelsiusTemp + 273
	return result, nil
}
