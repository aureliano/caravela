package provider

import "net/http"

type HTTPClientDecorator struct {
	Client http.Client
}

type HTTPClientPlugin interface {
	Do(req *http.Request) (*http.Response, error)
}

func (decorator *HTTPClientDecorator) Do(req *http.Request) (*http.Response, error) {
	return decorator.Client.Do(req)
}

// It is the interface that every release provider should implement, as it has all the
// expected method definitions for querying, caching and restoring cached release.
type UpdaterProvider interface {
	// FetchLastRelease queries provider for the last release of a project.
	FetchLastRelease(client HTTPClientPlugin) (*Release, error)

	// CacheRelease writes the release passed as parameter to the file system.
	CacheRelease(release Release) error

	// RestoreCacheRelease retores a cached release.
	// It returns nil if no release has been cached yet.
	RestoreCacheRelease() (*Release, error)
}
