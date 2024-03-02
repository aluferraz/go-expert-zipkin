package http_clients

import "net/http"

type ZipkinClientInterface interface {
	DoWithAppSpan(req *http.Request, span string) (*http.Response, error)
}
