package provider

import "net/http"

type HttpClientDecorator struct {
	Client http.Client
}

type HttpClientPlugin interface {
	Do(req *http.Request) (*http.Response, error)
}

func (decorator *HttpClientDecorator) Do(req *http.Request) (*http.Response, error) {
	return decorator.Client.Do(req)
}

type UpdaterProvider interface {
	FetchLastRelease(client HttpClientPlugin) (*Release, error)
	CacheRelease(release Release) error
	RestoreCacheRelease() (*Release, error)
}
