package mocks

import (
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
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
	c := &ZipkinMockClient{tracer: nil, Client: &http.Client{}}
	return c
}

func (c *ZipkinMockClient) DoWithAppSpan(req *http.Request, name string) (*http.Response, error) {
	return c.Client.Do(req)
}
