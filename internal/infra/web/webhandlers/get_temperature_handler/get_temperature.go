package get_temperature_handler

import (
	"encoding/json"
	"errors"
	zipcode2 "github.com/aluferraz/go-expert-zipkin/internal/entity/zipcode"
	"github.com/aluferraz/go-expert-zipkin/internal/infra/http_clients"
	"github.com/aluferraz/go-expert-zipkin/internal/usecase/get_temperature"
	"net/http"
)

type WebGetTemperatureHandler struct {
	usecase       get_temperature.UseCase
	client        http_clients.ZipkinClientInterface
	WeatherApiKey string
	WeatherApiUrl string
	ApiCepUrl     string
}

func NewGetTemperatureHandler(usecase get_temperature.UseCase, client http_clients.ZipkinClientInterface,
) *WebGetTemperatureHandler {
	return &WebGetTemperatureHandler{
		usecase: usecase,
		client:  client,
	}
}

func (h *WebGetTemperatureHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var dto get_temperature.InputDTO
	var err error

	zipcode_url := r.URL.Query().Get("zipcode")
	zipcode, err := zipcode2.NewZipcode(zipcode_url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dto = get_temperature.InputDTO{
		Zipcode: zipcode,
	}

	dto.WeatherApiKey = h.WeatherApiKey
	dto.WeatherApiUrl = h.WeatherApiUrl
	dto.ApiCepUrl = h.ApiCepUrl
	ctx := r.Context()
	dto.Ctx = &ctx

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	output, err := h.usecase.Execute(dto)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if output.City == "" {
		http.Error(w, errors.New("can not find zipcode").Error(), http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
