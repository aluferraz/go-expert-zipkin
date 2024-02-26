package get_temperature

import (
	"context"
	"encoding/json"
	zipcode2 "github.com/aluferraz/go-expert-zipkin/internal/entity/zipcode"
	"github.com/aluferraz/go-expert-zipkin/internal/infra/mocks"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type GetTemperatureTestSuite struct {
	suite.Suite
	ctx    context.Context
	client mocks.ZipkinClientInterface
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(GetTemperatureTestSuite))
}

func (suite *GetTemperatureTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.client = mocks.NewZipkinMockClient()
}

// Test case for getCityFromZipCode function using mock server
func (suite *GetTemperatureTestSuite) TestGetCityFromZipCode() {
	// Mock API response
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := AddressResponse{Localidade: "Dublin"}
		jsonBytes, _ := json.Marshal(response)
		_, _ = w.Write(jsonBytes)
	}))
	defer mockServer.Close()
	uc := NewUseCase(suite.client)

	// Call the function with the mock server URL
	city, err := uc.getCityFromZipCode(mockServer.URL+"/%s/json/", "12345", &suite.ctx)

	// Check if there's no error
	suite.NoError(err)
	suite.Equal("Dublin", city)
}

// Test case for getTemperatureFromCity function using mock server
func (suite *GetTemperatureTestSuite) TestGetTemperatureFromCity() {
	// Mock API response
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{"current": map[string]interface{}{"temp_c": 10.5}}
		jsonBytes, _ := json.Marshal(response)
		_, _ = w.Write(jsonBytes)
	}))
	defer mockServer.Close()

	// Call the function with the mock server URL
	uc := NewUseCase(suite.client)

	temperature, err := uc.getTemperatureFromCity(mockServer.URL+"/current.json?key=%s&q=%s", "your-api-key", "Dublin", &suite.ctx)

	// Check if there's no error
	suite.NoError(err)
	suite.Equal(10.5, temperature)
}
func (suite *GetTemperatureTestSuite) TestUseCase_Execute() {
	// Mock API responses
	mockApiCepServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := AddressResponse{Localidade: "Dublin"}
		jsonBytes, _ := json.Marshal(response)
		_, _ = w.Write(jsonBytes)
	}))
	defer mockApiCepServer.Close()

	mockWeatherApiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{"current": map[string]interface{}{"temp_c": 10.5}}
		jsonBytes, _ := json.Marshal(response)
		_, _ = w.Write(jsonBytes)
	}))
	defer mockWeatherApiServer.Close()

	// Create a test inputDTO with the URLs of the mock servers
	zipcode, err := zipcode2.NewZipcode("12345678")

	input := InputDTO{
		WeatherApiKey: "your-api-key",
		WeatherApiUrl: mockWeatherApiServer.URL + "/current.json?key=%s&q=%s",
		ApiCepUrl:     mockApiCepServer.URL + "/%s/json/",
		Zipcode:       zipcode, // replace with a valid zip code for testing
	}
	// Create the use case instance with test inputDTO, context and http client
	uc := NewUseCase(suite.client)

	// Execute the use case
	input.Ctx = &suite.ctx
	output, err := uc.Execute(input)

	// Check if there's no error
	suite.NoError(err)

	// Check if city is not empty
	suite.NotEmpty(output.City)

	// Check if temperatures are not zero
	suite.NotZero(output.CelsiusTemp)
	suite.NotZero(output.FahrenheitTemp)
	suite.NotZero(output.KelvinTemp)
}
