package updater

import (
	"fmt"

	pvdr "github.com/aureliano/caravela/provider"
)

// FindUpdate fetches the last release published.
//
// It returns the last release available or raises an error
// if the current version is already the last one.
func FindUpdate(
	client pvdr.HTTPClientPlugin,
	provider pvdr.UpdaterProvider,
	currver string,
) (*pvdr.Release, error) {
	rel, err := provider.RestoreCacheRelease()

	if err != nil {
		rel, err = provider.FetchLastRelease(client)
		if err != nil {
			return nil, err
		}

		_ = provider.CacheRelease(*rel)
	}

	if rel.CompareTo(&pvdr.Release{Name: currver}) == 1 {
		return rel, nil
	}

	return nil, fmt.Errorf("already on the edge")
}
