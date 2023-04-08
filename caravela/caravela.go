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
	I18nConf
	HttpClient *http.Client
}

var mpDownloadTo func(pvdr.HTTPClientPlugin, *pvdr.Release, string) (string, string, error) = downloadTo
var mpDecompress func(src string) (int, error) = decompress
var mpChecksum func(binPath string, checksumsPath string) error = checksum
var mpInstall func(srcDir string) error = install
var mpCheckForUpdates func(client pvdr.HTTPClientPlugin, provider pvdr.UpdaterProvider, currver string) (*pvdr.Release, error) = checkForUpdates
var mpUpdate func(client pvdr.HTTPClientPlugin, provider pvdr.UpdaterProvider, pname, currver string) error = update

// CheckForUpdates queries, given a provider, for new releases.
// It returns the last release available or nil if the current
// version is already the last one.
func CheckForUpdates(c Conf) (*pvdr.Release, error) {
	if c.Version == "" {
		return nil, fmt.Errorf("current version is required")
	}

	if c.HttpClient == nil {
		c.HttpClient = http.DefaultClient
	}

	client := pvdr.HTTPClientDecorator{Client: *c.HttpClient}

	err := prepareI18n(c.I18nConf)
	if err != nil {
		c.I18nConf = I18nConf{Verbose: false, Locale: EN}
		_ = prepareI18n(c.I18nConf)
		fmt.Println("Use default I18n configuration.")
	}

	return mpCheckForUpdates(&client, c.Provider, c.Version)
}

// Update running program to the last available release.
// Raises an error if it's already the last version.
func Update(c Conf) error {
	if c.ProcessName == "" {
		return fmt.Errorf("process name is required")
	}

	if c.HttpClient == nil {
		c.HttpClient = http.DefaultClient
	}

	client := pvdr.HTTPClientDecorator{Client: *c.HttpClient}

	err := prepareI18n(c.I18nConf)
	if err != nil {
		c.I18nConf = I18nConf{Verbose: false, Locale: EN}
		_ = prepareI18n(c.I18nConf)
		fmt.Println("Use default I18n configuration.")
	}

	return mpUpdate(&client, c.Provider, c.ProcessName, c.Version)
}

func checkForUpdates(client pvdr.HTTPClientPlugin, provider pvdr.UpdaterProvider, currver string) (*pvdr.Release, error) {
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
	} else {
		return nil, nil
	}
}

func update(client pvdr.HTTPClientPlugin, provider pvdr.UpdaterProvider, pname, currver string) error {
	rel, err := checkForUpdates(client, provider, currver)
	if err != nil {
		return err
	}

	wmsg(200)

	if rel == nil {
		return fmt.Errorf("already on the edge")
	}

	wmsg(201, rel.Name)
	fmt.Println(rel.Description)

	dir := filepath.Join(os.TempDir(), pname)
	_ = os.MkdirAll(dir, os.ModePerm)
	wmsg(202, dir)

	bin, checksums, err := mpDownloadTo(client, rel, dir)
	if err != nil {
		return err
	}

	wmsg(203)
	num, err := mpDecompress(bin)
	if err != nil {
		return err
	}
	wmsg(204, num, filepath.Base(bin))

	err = mpChecksum(bin, checksums)
	if err != nil {
		return err
	}

	err = mpInstall(dir)
	if err != nil {
		return err
	}

	wmsg(205)
	os.RemoveAll(dir)

	return nil
}
