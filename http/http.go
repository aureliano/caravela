package http

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

type HttpClientDecorator struct {
	Timeout    time.Duration
	TlsVersion uint16
	client     http.Client
}

type HttpClientPlugin interface {
	Do(req *http.Request) (*http.Response, error)
}

func (decorator *HttpClientDecorator) Do(req *http.Request) (*http.Response, error) {
	return decorator.client.Do(req)
}

func BuildClient(timeout uint32, tlsVersion int) (*HttpClientDecorator, error) {
	var version uint16
	switch tlsVersion {
	case 10:
		version = tls.VersionTLS10
	case 11:
		version = tls.VersionTLS11
	case 12:
		version = tls.VersionTLS12
	case 13:
		version = tls.VersionTLS13
	default:
		return nil, fmt.Errorf("unsupported tls version %d", tlsVersion)
	}

	client := http.Client{Timeout: time.Duration(timeout) * time.Second}
	client.Transport = &http.Transport{TLSClientConfig: &tls.Config{MinVersion: version}}

	return &HttpClientDecorator{Timeout: client.Timeout, TlsVersion: version, client: client}, nil
}

func BuildClientTls10() (*HttpClientDecorator, error) {
	return BuildClient(30, 10)
}

func BuildClientTls11() (*HttpClientDecorator, error) {
	return BuildClient(30, 11)
}

func BuildClientTls12() (*HttpClientDecorator, error) {
	return BuildClient(30, 12)
}

func BuildClientTls13() (*HttpClientDecorator, error) {
	return BuildClient(30, 13)
}
