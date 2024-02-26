package mocks

import "net/http"

type ZipkinClientInterface interface {
	DoWithAppSpan(req *http.Request, name string) (*http.Response, error)
}
