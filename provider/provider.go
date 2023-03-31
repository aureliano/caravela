package provider

import (
	"github.com/aureliano/caravela/http"
	"github.com/aureliano/caravela/release"
)

type UpdaterProvider interface {
	FetchReleases(client http.HttpClientPlugin) ([]*release.Release, error)
	FetchLastRelease(client http.HttpClientPlugin) (*release.Release, error)
	CacheRelease(release release.Release) error
	RestoreCacheRelease() (*release.Release, error)
}
