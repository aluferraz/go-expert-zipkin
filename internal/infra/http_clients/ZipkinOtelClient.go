package http_clients

import (
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
)

type ZipkinMockClient struct {
	*http.Client
	tracer           *zipkin.Tracer
	httpTrace        bool
	defaultTags      map[string]string
	transportOptions []http.Transport
	remoteEndpoint   *model.Endpoint
}

func NewZipkinMockClient() *ZipkinMockClient {
	c := &ZipkinMockClient{tracer: nil, Client: otelhttp.DefaultClient}
	return c
}

func (c *ZipkinMockClient) DoWithAppSpan(req *http.Request, name string) (*http.Response, error) {
	ctx := req.Context()
	tracer := otel.Tracer("zipkin-goexpert")
	carrier := propagation.HeaderCarrier(req.Header)
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	ctx, httpSpan := tracer.Start(ctx, name)
	defer httpSpan.End()
	return c.Client.Do(req)
}
