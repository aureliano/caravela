package caravela

import (
	"github.com/aureliano/caravela/http"
)

type UpdaterProvider interface {
	FetchLastRelease(client http.HttpClientPlugin) (*Release, error)
	CacheRelease(release Release) error
	RestoreCacheRelease() (*Release, error)
}
