package temperature_input

import (
	"bytes"
	"context"
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

type TemperatureInputTestSuite struct {
	suite.Suite
	ctx context.Context
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TemperatureInputTestSuite))
}

func (suite *TemperatureInputTestSuite) SetupSuite() {
}
func (s *TemperatureInputTestSuite) TestTemperatureInput() {
	clientMock := mocks.NewZipkinMockClient()
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	// Mock API responses

	mockWeatherApiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{"city": "São Paulo", "temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.65}
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
	handler := NewTemperatureInputHandler(mockWeatherApiServer.URL, clientMock)

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
