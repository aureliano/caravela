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

// UpdaterProvider is the interface that all providers must implement.
// It has the basic methods which any provider must have implemented,
// that is expected on the core package.
type UpdaterProvider interface {
	// FetchLastRelease queries provider for the last release of a project.
	FetchLastRelease(client HttpClientPlugin) (*Release, error)

	// CacheRelease caches the release passed as parameter on file system.
	CacheRelease(release Release) error

	// RestoreCacheRelease retores a cached release.
	// It returns nil if any release wasn't cached yet.
	RestoreCacheRelease() (*Release, error)
}
