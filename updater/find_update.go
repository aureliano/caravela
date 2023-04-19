package updater

import (
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
	ignoreCache bool,
) (*pvdr.Release, error) {
	var release *pvdr.Release
	var err error

	if ignoreCache {
		release, err = provider.FetchLastRelease(client)
	} else {
		release, err = findUpdateUseCache(client, provider)
	}

	if err != nil {
		return nil, err
	}

	if release.CompareTo(&pvdr.Release{Name: currver}) == 1 {
		return release, nil
	}

	return &pvdr.Release{}, nil
}

func findUpdateUseCache(client pvdr.HTTPClientPlugin, provider pvdr.UpdaterProvider) (*pvdr.Release, error) {
	release, err := provider.RestoreCacheRelease()

	if err != nil {
		release, err = provider.FetchLastRelease(client)
		if err != nil {
			return nil, err
		}

		_ = provider.CacheRelease(*release)
	}

	return release, nil
}
