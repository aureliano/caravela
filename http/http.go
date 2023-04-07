package http

import (
	"net/http"
)

type HttpClientDecorator struct {
	client http.Client
}

type HttpClientPlugin interface {
	Do(req *http.Request) (*http.Response, error)
}

func (decorator *HttpClientDecorator) Do(req *http.Request) (*http.Response, error) {
	return decorator.client.Do(req)
}
