package caravela

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	pvdr "github.com/aureliano/caravela/provider"
)

// A Conf is a wrapper o data to be passed as input to the public functions.
type Conf struct {
	ProcessName string
	Version     string
	Provider    pvdr.UpdaterProvider
	HTTPClient  *http.Client
}

var mpDownloadTo = downloadTo
var mpDecompress = decompress
var mpChecksum = checksum
var mpInstall = install
var mpCheckForUpdates = checkForUpdates
var mpUpdate = update

// CheckForUpdates queries, given a provider, for new releases.
// It returns the last release available or nil if the current
// version is already the last one.
func CheckForUpdates(c Conf) (*pvdr.Release, error) {
	if c.Version == "" {
		return nil, fmt.Errorf("current version is required")
	}

	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}

	client := pvdr.HTTPClientDecorator{Client: *c.HTTPClient}

	return mpCheckForUpdates(&client, c.Provider, c.Version)
}

// Update running program to the last available release.
// Raises an error if it's already the last version
// or returns the new release.
func Update(c Conf) (*pvdr.Release, error) {
	if c.ProcessName == "" {
		return nil, fmt.Errorf("process name is required")
	}

	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}

	client := pvdr.HTTPClientDecorator{Client: *c.HTTPClient}

	return mpUpdate(&client, c.Provider, c.ProcessName, c.Version)
}

func checkForUpdates(client pvdr.HTTPClientPlugin, provider pvdr.UpdaterProvider,
	currver string) (*pvdr.Release, error) {
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

func update(client pvdr.HTTPClientPlugin, provider pvdr.UpdaterProvider, pname, currver string) (*pvdr.Release, error) {
	rel, err := checkForUpdates(client, provider, currver)
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(os.TempDir(), pname)
	_ = os.MkdirAll(dir, os.ModePerm)

	bin, checksums, err := mpDownloadTo(client, rel, dir)
	if err != nil {
		return nil, err
	}

	_, err = mpDecompress(bin)
	if err != nil {
		return nil, err
	}

	err = mpChecksum(bin, checksums)
	if err != nil {
		return nil, err
	}

	err = mpInstall(dir)
	if err != nil {
		return nil, err
	}

	os.RemoveAll(dir)

	return rel, nil
}
