package caravela

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	pvdr "github.com/aureliano/caravela/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockHTTPClient struct{ mock.Mock }

type mockProvider struct{ mock.Mock }

func (client *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := client.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (provider *mockProvider) FetchReleases(client pvdr.HTTPClientPlugin) ([]*pvdr.Release, error) {
	args := provider.Called(client)
	return args.Get(0).([]*pvdr.Release), args.Error(1)
}

func (provider *mockProvider) FetchLastRelease(client pvdr.HTTPClientPlugin) (*pvdr.Release, error) {
	args := provider.Called(client)
	var rel *pvdr.Release
	if args.Get(0) != nil {
		rel, _ = args.Get(0).(*pvdr.Release)
	}

	return rel, args.Error(1)
}

func (provider *mockProvider) CacheRelease(rel pvdr.Release) error {
	args := provider.Called(rel)
	return args.Error(0)
}

func (provider *mockProvider) RestoreCacheRelease() (*pvdr.Release, error) {
	args := provider.Called()
	var rel *pvdr.Release
	if args.Get(0) != nil {
		rel, _ = args.Get(0).(*pvdr.Release)
	}

	return rel, args.Error(1)
}

func TestCheckForUpdatesCurrentVersionIsRequired(t *testing.T) {
	_, err := CheckForUpdates(Conf{})
	assert.Equal(t, "current version is required", err.Error())
}

func TestCheckForUpdatesI18nError(t *testing.T) {
	mpCheckForUpdates = func(client pvdr.HTTPClientPlugin, provider pvdr.UpdaterProvider,
		currver string) (*pvdr.Release, error) {
		return nil, fmt.Errorf("already on the edge")
	}

	r, err := CheckForUpdates(Conf{I18nConf: I18nConf{Verbose: false, Locale: -1}, Version: "0.1.0"})
	assert.Nil(t, r)
	assert.Equal(t, "already on the edge", err.Error())
}

func TestCheckForUpdates(t *testing.T) {
	mpCheckForUpdates = func(client pvdr.HTTPClientPlugin, provider pvdr.UpdaterProvider,
		currver string) (*pvdr.Release, error) {
		return nil, fmt.Errorf("already on the edge")
	}

	r, err := CheckForUpdates(Conf{I18nConf: I18nConf{Verbose: false, Locale: PtBr},
		HTTPClient: http.DefaultClient, Version: "0.1.0"})
	assert.Nil(t, r)
	assert.Equal(t, "already on the edge", err.Error())
}

func TestUpdateProcessNameIsRequired(t *testing.T) {
	err := Update(Conf{ProcessName: "", Version: "0.1.0"})
	assert.Equal(t, "process name is required", err.Error())
}

func TestUpdateI18nError(t *testing.T) {
	mpUpdate = func(client pvdr.HTTPClientPlugin, provider pvdr.UpdaterProvider, pname, currver string) error {
		return nil
	}

	err := Update(Conf{I18nConf: I18nConf{Verbose: false, Locale: -1}, ProcessName: "oalienista", Version: "0.1.0"})
	assert.Nil(t, err)
}

func TestCheckForUpdatesRestoreCacheCurrentVersionIsOlder(t *testing.T) {
	m := new(mockHTTPClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProvider)
	p.On("RestoreCacheRelease").Return(&pvdr.Release{Name: "v0.1.0"}, nil)

	r, _ := checkForUpdates(m, p, "v0.1.0-alpha")
	assert.Equal(t, r.Name, "v0.1.0")
	p.AssertCalled(t, "RestoreCacheRelease")
}

func TestCheckForUpdatesRestoreCacheCurrentVersionOnTheEdge(t *testing.T) {
	m := new(mockHTTPClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProvider)
	p.On("RestoreCacheRelease").Return(&pvdr.Release{Name: "v0.1.0"}, nil)

	r, _ := checkForUpdates(m, p, "v0.1.0")
	assert.Nil(t, r)
	p.AssertCalled(t, "RestoreCacheRelease")
}

func TestCheckForUpdatesNoCacheFetchLastReleaseError(t *testing.T) {
	m := new(mockHTTPClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProvider)
	p.On("CacheRelease", pvdr.Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("any error"))
	p.On("FetchLastRelease", m).Return(
		nil, fmt.Errorf("some error"),
	)

	r, e := checkForUpdates(m, p, "v0.1.2")
	assert.Nil(t, r)
	assert.Equal(t, "some error", e.Error())
	p.AssertCalled(t, "FetchLastRelease", m)
	p.AssertCalled(t, "RestoreCacheRelease")
}

func TestCheckForUpdatesNoCacheCurrentVersionIsOlder(t *testing.T) {
	m := new(mockHTTPClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProvider)
	p.On("CacheRelease", pvdr.Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("any error"))
	p.On("FetchLastRelease", m).Return(
		&pvdr.Release{
			Name: "v0.1.2",
		}, nil,
	)

	r, _ := checkForUpdates(m, p, "v0.1.1")
	assert.Equal(t, r.Name, "v0.1.2")
	p.AssertCalled(t, "FetchLastRelease", m)
	p.AssertCalled(t, "CacheRelease", pvdr.Release{Name: "v0.1.2"})
	p.AssertCalled(t, "RestoreCacheRelease")
}

func TestCheckForUpdatesNoCacheCurrentVersionOnTheEdge(t *testing.T) {
	m := new(mockHTTPClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProvider)
	p.On("CacheRelease", pvdr.Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("any error"))
	p.On("FetchLastRelease", m).Return(
		&pvdr.Release{
			Name: "v0.1.2",
		}, nil,
	)

	r, _ := checkForUpdates(m, p, "v0.1.2")
	assert.Nil(t, r)
	p.AssertCalled(t, "FetchLastRelease", m)
	p.AssertCalled(t, "CacheRelease", pvdr.Release{Name: "v0.1.2"})
	p.AssertCalled(t, "RestoreCacheRelease")
}

func TestUpdateCheckVersionFail(t *testing.T) {
	m := new(mockHTTPClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(``))),
		}, nil)
	p := new(mockProvider)
	p.On("FetchLastRelease", m).Return(
		nil, fmt.Errorf("any error"),
	)
	p.On("CacheRelease", pvdr.Release{}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("any error"))

	err := update(m, p, "14-bis", "0.0.1")
	actual := err.Error()
	expected := "any error"

	p.AssertNotCalled(t, "FetchLastRelease")
	p.AssertNotCalled(t, "CacheRelease")
	p.AssertCalled(t, "RestoreCacheRelease")
	assert.Equal(t, expected, actual)
}

func TestUpdateAlreadyUpToDate(t *testing.T) {
	m := new(mockHTTPClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProvider)
	p.On("FetchLastRelease", m).Return(
		&pvdr.Release{
			Name: "v0.1.2",
		}, nil,
	)
	p.On("CacheRelease", pvdr.Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("no file error"))

	err := update(m, p, "14-bis", "0.1.2")
	actual := err.Error()
	expected := "already on the edge"

	p.AssertNotCalled(t, "FetchLastRelease")
	p.AssertNotCalled(t, "CacheRelease")
	p.AssertCalled(t, "RestoreCacheRelease")
	assert.Equal(t, expected, actual)
}

func TestUpdateDownloadFail(t *testing.T) {
	m := new(mockHTTPClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProvider)
	p.On("FetchLastRelease", m).Return(
		&pvdr.Release{
			Name: "v0.1.2",
		}, nil,
	)
	p.On("CacheRelease", pvdr.Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("no file error"))
	mpDownloadTo = func(hcp pvdr.HTTPClientPlugin, r *pvdr.Release, s string) (string, string, error) {
		return "", "", fmt.Errorf("download release error")
	}

	err := update(m, p, "14-bis", "0.1.1")

	p.AssertNotCalled(t, "FetchLastRelease")
	p.AssertNotCalled(t, "CacheRelease")
	p.AssertCalled(t, "RestoreCacheRelease")
	assert.Equal(t, "download release error", err.Error())
}

func TestUpdateDecompressionFail(t *testing.T) {
	m := new(mockHTTPClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProvider)
	p.On("FetchLastRelease", m).Return(
		&pvdr.Release{
			Name: "v0.1.2",
		}, nil,
	)
	p.On("CacheRelease", pvdr.Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("no file error"))
	mpDownloadTo = func(hcp pvdr.HTTPClientPlugin, r *pvdr.Release, s string) (string, string, error) {
		return "", "", nil
	}
	mpDecompress = func(src string) (int, error) { return 0, fmt.Errorf("decompression error") }
	err := update(m, p, "14-bis", "0.1.1")

	p.AssertNotCalled(t, "FetchLastRelease")
	p.AssertNotCalled(t, "CacheRelease")
	p.AssertCalled(t, "RestoreCacheRelease")
	assert.Equal(t, "decompression error", err.Error())
}

func TestUpdateChecksumFail(t *testing.T) {
	m := new(mockHTTPClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProvider)
	p.On("FetchLastRelease", m).Return(
		&pvdr.Release{
			Name: "v0.1.2",
		}, nil,
	)
	p.On("CacheRelease", pvdr.Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("no file error"))
	mpDownloadTo = func(hcp pvdr.HTTPClientPlugin, r *pvdr.Release, s string) (string, string, error) {
		return "", "", nil
	}
	mpDecompress = func(src string) (int, error) { return 1, nil }
	mpChecksum = func(binPath, checksumsPath string) error { return fmt.Errorf("checksum error") }
	err := update(m, p, "14-bis", "0.1.1")

	p.AssertNotCalled(t, "FetchLastRelease")
	p.AssertNotCalled(t, "CacheRelease")
	p.AssertCalled(t, "RestoreCacheRelease")
	assert.Equal(t, "checksum error", err.Error())
}

func TestUpdateInstallationFail(t *testing.T) {
	m := new(mockHTTPClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProvider)
	p.On("FetchLastRelease", m).Return(
		&pvdr.Release{
			Name: "v0.1.2",
		}, nil,
	)
	p.On("CacheRelease", pvdr.Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("no file error"))
	mpDownloadTo = func(hcp pvdr.HTTPClientPlugin, r *pvdr.Release, s string) (string, string, error) {
		return "", "", nil
	}
	mpDecompress = func(src string) (int, error) { return 1, nil }
	mpChecksum = func(binPath, checksumsPath string) error { return nil }
	mpInstall = func(srcDir string) error { return fmt.Errorf("installation error") }
	err := update(m, p, "14-bis", "0.1.1")

	p.AssertNotCalled(t, "FetchLastRelease")
	p.AssertNotCalled(t, "CacheRelease")
	p.AssertCalled(t, "RestoreCacheRelease")
	assert.Equal(t, "installation error", err.Error())
}

func TestUpdate(t *testing.T) {
	m := new(mockHTTPClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProvider)
	p.On("FetchLastRelease", m).Return(
		&pvdr.Release{
			Name: "v0.1.2",
		}, nil,
	)
	p.On("CacheRelease", pvdr.Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("no file error"))
	mpDownloadTo = func(hcp pvdr.HTTPClientPlugin, r *pvdr.Release, s string) (string, string, error) {
		return "", "", nil
	}
	mpDecompress = func(src string) (int, error) { return 1, nil }
	mpChecksum = func(binPath, checksumsPath string) error { return nil }
	mpInstall = func(srcDir string) error { return nil }
	err := update(m, p, "14-bis", "0.1.1")

	p.AssertNotCalled(t, "FetchLastRelease")
	p.AssertNotCalled(t, "CacheRelease")
	p.AssertCalled(t, "RestoreCacheRelease")
	assert.Nil(t, err)
}
