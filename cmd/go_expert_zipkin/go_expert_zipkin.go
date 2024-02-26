package go_expert_zipkin

import (
	"context"
	"fmt"
	"github.com/aluferraz/go-expert-zipkin/cmd/go_expert_zipkin/dependency_injection"
	"github.com/aluferraz/go-expert-zipkin/configs"
	"github.com/aluferraz/go-expert-zipkin/internal/infra/web/webhandlers/temperature_input"
	"github.com/aluferraz/go-expert-zipkin/internal/infra/web/webserver"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	reporterhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"log"
	"net/http"
	"os"
)

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Bootstap() {
	workdir, err := os.Getwd()
	handleErr(err)
	appConfig, err := configs.LoadConfig(workdir)
	if err != nil {
		panic(err)
	}

	//Zipkin
	const endpointURL = "http://observability:9411/api/v2/spans"

	reporter := reporterhttp.NewReporter(endpointURL)
	defer func() {
		_ = reporter.Close()
	}()
	// create our local service endpoint
	endpoint, err := zipkin.NewEndpoint("go-expert-zipkin", "localhost:0")
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	// initialize our tracer
	tracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}
	// create global zipkin traced http client
	client, err := zipkinhttp.NewClient(tracer, zipkinhttp.ClientTrace(true))
	if err != nil {
		log.Fatalf("unable to create client: %+v\n", err)
	}

	restServer := webserver.NewWebServer(appConfig.WebserverPort, tracer)
	ctx := context.Background()

	temperatureHandler := dependency_injection.NewTemperatureHandler(&ctx, client)
	temperatureHandler.WeatherApiKey = appConfig.WeatherApiKey
	temperatureHandler.ApiCepUrl = appConfig.CepApiURL
	temperatureHandler.WeatherApiUrl = appConfig.WeatherApiURL

	/*restServer.AddHandler("/", http.MethodGet, temperatureHandler.Handle)*/
	restServer.AddHandler("/servicoB", http.MethodGet, temperatureHandler.Handle)

	temperatureInputHandler := temperature_input.NewTemperatureInputHandler(
		fmt.Sprintf("http://localhost%s/servicoB", restServer.WebServerPort),
		client,
	)

	restServer.AddHandler("/", http.MethodGet, temperatureInputHandler.Handle)
	restServer.Start()

}
