package caravela

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type Conf struct {
	ProcessName string
	Version     string
	Provider    UpdaterProvider
	I18nConf
	HttpClient http.Client
}

type HttpClientDecorator struct {
	client http.Client
}

type HttpClientPlugin interface {
	Do(req *http.Request) (*http.Response, error)
}

func (decorator *HttpClientDecorator) Do(req *http.Request) (*http.Response, error) {
	return decorator.client.Do(req)
}

var downloadRelease func(HttpClientPlugin, *Release, string) (string, string, error) = downloadTo
var funcDecompress func(src string) (int, error) = decompress
var checksumRelease func(binPath string, checksumsPath string) error = checksum
var installRelease func(srcDir string) error = install

func CheckForUpdates(client HttpClientPlugin, provider UpdaterProvider, conf I18nConf, currver string) (*Release, error) {
	err := prepareI18n(conf)
	if err != nil {
		_ = prepareI18n(I18nConf{Verbose: false, Locale: EN})
		fmt.Println("Use default I18n configuration.")
	}

	rel, err := provider.RestoreCacheRelease()

	if err != nil {
		rel, err = provider.FetchLastRelease(client)
		if err != nil {
			return nil, err
		}

		_ = provider.CacheRelease(*rel)
	}

	if rel.CompareTo(&Release{Name: currver}) == 1 {
		return rel, nil
	} else {
		return nil, nil
	}
}

func Update(client HttpClientPlugin, provider UpdaterProvider, conf I18nConf, pname, currver string) error {
	rel, err := CheckForUpdates(client, provider, conf, currver)
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

	bin, checksums, err := downloadRelease(client, rel, dir)
	if err != nil {
		return err
	}

	wmsg(203)
	num, err := funcDecompress(bin)
	if err != nil {
		return err
	}
	wmsg(204, num, filepath.Base(bin))

	err = checksumRelease(bin, checksums)
	if err != nil {
		return err
	}

	err = installRelease(dir)
	if err != nil {
		return err
	}

	wmsg(205)
	os.RemoveAll(dir)

	return nil
}
