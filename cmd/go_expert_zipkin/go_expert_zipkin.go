package go_expert_zipkin

import (
	"context"
	"fmt"
	"github.com/aluferraz/go-expert-zipkin/cmd/go_expert_zipkin/dependency_injection"
	"github.com/aluferraz/go-expert-zipkin/configs"
	"github.com/aluferraz/go-expert-zipkin/internal/infra/http_clients"
	"github.com/aluferraz/go-expert-zipkin/internal/infra/web/webhandlers/temperature_input"
	"github.com/aluferraz/go-expert-zipkin/internal/infra/web/webserver"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"google.golang.org/grpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	//"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"net/http"
	"os"
)

func handleErr(err error) {
	if err != nil {
		log.Fatal().Err(err)
	}
}

const endpointURL = "otel-collector:4317"

func Bootstap() {
	workdir, err := os.Getwd()
	handleErr(err)
	appConfig, err := configs.LoadConfig(workdir)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	// Initialize OpenTelemetry Tracer Provider
	shutdown := initTracer(ctx)
	defer shutdown(context.Background())

	// create global zipkin traced http client
	ctx, span := otel.Tracer("zipkin-goexpert").Start(ctx, "main-handler")
	defer span.End()
	client := http_clients.NewZipkinMockClient()

	restServer := webserver.NewWebServer(appConfig.WebserverPort)

	temperatureHandler := dependency_injection.NewTemperatureHandler(&ctx, client)
	temperatureHandler.WeatherApiKey = appConfig.WeatherApiKey
	temperatureHandler.ApiCepUrl = appConfig.CepApiURL
	temperatureHandler.WeatherApiUrl = appConfig.WeatherApiURL

	/*restServer.AddHandler("/", http.MethodGet, temperatureHandler.Handle)*/
	restServer.AddHandler("/servicoB", http.MethodGet, temperatureHandler.Handle, "servicoB")

	temperatureInputHandler := temperature_input.NewTemperatureInputHandler(
		fmt.Sprintf("http://web_a%s/servicoB", restServer.WebServerPort),
		client,
	)

	restServer.AddHandler("/", http.MethodPost, temperatureInputHandler.Handle, "servicoA")
	restServer.Start()

}
func initTracer(ctx context.Context) func( context.Context) error {
	//exporter, err := zipkin.New(
	//	endpointURL,
	//	// Additional Zipkin exporter options if desired
	//)
	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpointURL),
		otlptracegrpc.WithDialOption(grpc.WithBlock()))
	exporter, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		log.Fatal().Msgf("failed to create Otel exporter: %v", err)
	}
	bsp := sdktrace.NewBatchSpanProcessor(exporter)

	// Create a resource with service details
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("zipkin-goexpert"), // Your service name
			// Add other attributes as needed
		),
	)
	if err != nil {
		log.Fatal().Msgf("failed to create resource: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tp)
	propagator := b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader | b3.B3SingleHeader))

	// Set global propagator for context propagation
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagator,
		propagation.TraceContext{},
		propagation.Baggage{}),
	)

	return tp.Shutdown
}
