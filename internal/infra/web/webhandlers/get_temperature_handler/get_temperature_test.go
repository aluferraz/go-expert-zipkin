package get_temperature_handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	zipcode2 "github.com/aluferraz/go-expert-zipkin/internal/entity/zipcode"
	"github.com/aluferraz/go-expert-zipkin/internal/infra/mocks"
	"github.com/aluferraz/go-expert-zipkin/internal/usecase/get_temperature"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type GetTemperatureTestSuite struct {
	suite.Suite
}

func TestGetTemperatureTestSuite(t *testing.T) {
	suite.Run(t, new(GetTemperatureTestSuite))
}

func (s *GetTemperatureTestSuite) SetupTest() {}

func (s *GetTemperatureTestSuite) TestWebGetTemperatureHandler_Handle() {
	clientMock := mocks.NewZipkinMockClient()

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Mock API responses
	mockApiCepServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := get_temperature.AddressResponse{Localidade: "Dublin"}
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

	zipcode, err := zipcode2.NewZipcode("12345678")
	s.NoError(err)

	// Create a test inputDTO with the URLs of the mock servers
	input := get_temperature.InputDTO{
		Zipcode: zipcode,
	}
	uc := get_temperature.NewUseCase(clientMock)
	handler := NewGetTemperatureHandler(uc, clientMock)
	handler.WeatherApiUrl = mockWeatherApiServer.URL + "/current.json?key=%s&q=%s"
	handler.ApiCepUrl = mockApiCepServer.URL + "/%s/json/"

	jsonStr, err := json.Marshal(input)
	s.NoError(err)
	req, err := http.NewRequest("GET", fmt.Sprintf("/?zipcode=%s", zipcode.Zipcode), bytes.NewBuffer(jsonStr))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.Handle(rr, req)

	// Check the status code is what we expect.
	s.Equal(http.StatusOK, rr.Code)

	var output get_temperature.OutputDTO
	err = json.Unmarshal(rr.Body.Bytes(), &output)
	s.NoError(err)

	// Check if temperatures are not zero
	s.NotZero(output.CelsiusTemp)
	s.NotZero(output.FahrenheitTemp)
	s.NotZero(output.KelvinTemp)
}
