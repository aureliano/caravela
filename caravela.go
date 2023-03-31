package caravela

import (
	"github.com/aureliano/caravela/http"
	"github.com/aureliano/caravela/provider"
	"github.com/aureliano/caravela/release"
)

func CheckForUpdates(client http.HttpClientPlugin, provider provider.UpdaterProvider, currver string) (*release.Release, error) {
	rel, err := provider.RestoreCacheRelease()

	if err != nil {
		rel, err = provider.FetchLastRelease(client)
		if err != nil {
			return nil, err
		}

		_ = provider.CacheRelease(*rel)
	}

	if rel.CompareTo(&release.Release{Name: currver}) == 1 {
		return rel, nil
	} else {
		return nil, nil
	}
}
