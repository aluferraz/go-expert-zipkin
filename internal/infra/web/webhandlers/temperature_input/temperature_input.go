package temperature_input

import (
	"fmt"
	zipcode2 "github.com/aluferraz/go-expert-zipkin/internal/entity/zipcode"
	"github.com/aluferraz/go-expert-zipkin/internal/infra/http_clients"
	"github.com/aluferraz/go-expert-zipkin/internal/usecase/get_temperature"
	"io"
	"net/http"
	"net/url"
	"encoding/json"
)

type WebTemperatureInputHandler struct {
	service2Url string
	client      http_clients.ZipkinClientInterface
}
type InputDTO struct {
	Zipcode zipcode2.Zipcode `json:"cep"`
}

type FreeTextInput struct {
	Zipcode string `json:"cep"`
}

func NewTemperatureInputHandler(
	service2Url string,
	client http_clients.ZipkinClientInterface,
) *WebTemperatureInputHandler {
	return &WebTemperatureInputHandler{
		service2Url: service2Url,
		client:      client,
	}
}

func (h *WebTemperatureInputHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var dto get_temperature.InputDTO
	var user_input FreeTextInput
	var err error
	// zipcode_url := r.URL.Query().Get("zipcode")
	err = json.NewDecoder(r.Body).Decode(&user_input)
	zipcode, err := zipcode2.NewZipcode(user_input.Zipcode)
	if err != nil {
		/*http.Error(w, err.Error(), http.StatusBadRequest)*/
		http.Error(w, "invalid zipcode", http.StatusBadRequest)
		return
	}
	dto = get_temperature.InputDTO{
		Zipcode: zipcode,
	}

	url := fmt.Sprintf("%s?zipcode=%s", h.service2Url, url.QueryEscape(dto.Zipcode.Zipcode))
	newRequest, err := http.NewRequestWithContext(r.Context(), "GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := h.client.DoWithAppSpan(newRequest, "servico_b")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		Body.Close()
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	_, err = w.Write(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}