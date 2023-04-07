package caravela

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockHttpClient struct{ mock.Mock }

type mockProvider struct{ mock.Mock }

func (client *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := client.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (provider *mockProvider) FetchReleases(client HttpClientPlugin) ([]*Release, error) {
	args := provider.Called(client)
	return args.Get(0).([]*Release), args.Error(1)
}

func (provider *mockProvider) FetchLastRelease(client HttpClientPlugin) (*Release, error) {
	args := provider.Called(client)
	var rel *Release
	if args.Get(0) != nil {
		rel = args.Get(0).(*Release)
	}

	return rel, args.Error(1)
}

func (provider *mockProvider) CacheRelease(rel Release) error {
	args := provider.Called(rel)
	return args.Error(0)
}

func (provider *mockProvider) RestoreCacheRelease() (*Release, error) {
	args := provider.Called()
	var rel *Release
	if args.Get(0) != nil {
		rel = args.Get(0).(*Release)
	}

	return rel, args.Error(1)
}

func TestDo(t *testing.T) {
	c := HttpClientDecorator{client: http.Client{}}
	res, err := c.Do(&http.Request{})
	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestCheckForUpdatesRestoreCacheCurrentVersionIsOlder(t *testing.T) {
	m := new(mockHttpClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProvider)
	p.On("RestoreCacheRelease").Return(&Release{Name: "v0.1.0"}, nil)

	r, _ := CheckForUpdates(m, p, I18nConf{Verbose: true, Locale: -1}, "v0.1.0-alpha")
	assert.Equal(t, r.Name, "v0.1.0")
	p.AssertCalled(t, "RestoreCacheRelease")
}

func TestCheckForUpdatesRestoreCacheCurrentVersionOnTheEdge(t *testing.T) {
	m := new(mockHttpClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProvider)
	p.On("RestoreCacheRelease").Return(&Release{Name: "v0.1.0"}, nil)

	r, _ := CheckForUpdates(m, p, I18nConf{Verbose: true, Locale: EN}, "v0.1.0")
	assert.Nil(t, r)
	p.AssertCalled(t, "RestoreCacheRelease")
}

func TestCheckForUpdatesNoCacheFetchLastReleaseError(t *testing.T) {
	m := new(mockHttpClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProvider)
	p.On("CacheRelease", Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("any error"))
	p.On("FetchLastRelease", m).Return(
		nil, fmt.Errorf("some error"),
	)

	r, e := CheckForUpdates(m, p, I18nConf{Verbose: true, Locale: EN}, "v0.1.2")
	assert.Nil(t, r)
	assert.Equal(t, "some error", e.Error())
	p.AssertCalled(t, "FetchLastRelease", m)
	p.AssertCalled(t, "RestoreCacheRelease")
}

func TestCheckForUpdatesNoCacheCurrentVersionIsOlder(t *testing.T) {
	m := new(mockHttpClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProvider)
	p.On("CacheRelease", Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("any error"))
	p.On("FetchLastRelease", m).Return(
		&Release{
			Name: "v0.1.2",
		}, nil,
	)

	r, _ := CheckForUpdates(m, p, I18nConf{Verbose: true, Locale: EN}, "v0.1.1")
	assert.Equal(t, r.Name, "v0.1.2")
	p.AssertCalled(t, "FetchLastRelease", m)
	p.AssertCalled(t, "CacheRelease", Release{Name: "v0.1.2"})
	p.AssertCalled(t, "RestoreCacheRelease")
}

func TestCheckForUpdatesNoCacheCurrentVersionOnTheEdge(t *testing.T) {
	m := new(mockHttpClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProvider)
	p.On("CacheRelease", Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("any error"))
	p.On("FetchLastRelease", m).Return(
		&Release{
			Name: "v0.1.2",
		}, nil,
	)

	r, _ := CheckForUpdates(m, p, I18nConf{Verbose: true, Locale: EN}, "v0.1.2")
	assert.Nil(t, r)
	p.AssertCalled(t, "FetchLastRelease", m)
	p.AssertCalled(t, "CacheRelease", Release{Name: "v0.1.2"})
	p.AssertCalled(t, "RestoreCacheRelease")
}

func TestUpdateCheckVersionFail(t *testing.T) {
	m := new(mockHttpClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(``))),
		}, nil)
	p := new(mockProvider)
	p.On("FetchLastRelease", m).Return(
		nil, fmt.Errorf("any error"),
	)
	p.On("CacheRelease", Release{}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("any error"))

	err := Update(m, p, I18nConf{Verbose: true, Locale: EN}, "14-bis", "0.0.1")
	actual := err.Error()
	expected := "any error"

	p.AssertNotCalled(t, "FetchLastRelease")
	p.AssertNotCalled(t, "CacheRelease")
	p.AssertCalled(t, "RestoreCacheRelease")
	assert.Equal(t, expected, actual)
}

func TestUpdateAlreadyUpToDate(t *testing.T) {
	m := new(mockHttpClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProvider)
	p.On("FetchLastRelease", m).Return(
		&Release{
			Name: "v0.1.2",
		}, nil,
	)
	p.On("CacheRelease", Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("no file error"))

	err := Update(m, p, I18nConf{Verbose: true, Locale: EN}, "14-bis", "0.1.2")
	actual := err.Error()
	expected := "already on the edge"

	p.AssertNotCalled(t, "FetchLastRelease")
	p.AssertNotCalled(t, "CacheRelease")
	p.AssertCalled(t, "RestoreCacheRelease")
	assert.Equal(t, expected, actual)
}

func TestUpdateDownloadFail(t *testing.T) {
	m := new(mockHttpClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProvider)
	p.On("FetchLastRelease", m).Return(
		&Release{
			Name: "v0.1.2",
		}, nil,
	)
	p.On("CacheRelease", Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("no file error"))
	downloadRelease = func(hcp HttpClientPlugin, r *Release, s string) (string, string, error) {
		return "", "", fmt.Errorf("download release error")
	}

	err := Update(m, p, I18nConf{Verbose: true, Locale: EN}, "14-bis", "0.1.1")

	p.AssertNotCalled(t, "FetchLastRelease")
	p.AssertNotCalled(t, "CacheRelease")
	p.AssertCalled(t, "RestoreCacheRelease")
	assert.Equal(t, "download release error", err.Error())
}

func TestUpdateDecompressionFail(t *testing.T) {
	m := new(mockHttpClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProvider)
	p.On("FetchLastRelease", m).Return(
		&Release{
			Name: "v0.1.2",
		}, nil,
	)
	p.On("CacheRelease", Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("no file error"))
	downloadRelease = func(hcp HttpClientPlugin, r *Release, s string) (string, string, error) {
		return "", "", nil
	}
	funcDecompress = func(src string) (int, error) { return 0, fmt.Errorf("decompression error") }
	err := Update(m, p, I18nConf{Verbose: true, Locale: EN}, "14-bis", "0.1.1")

	p.AssertNotCalled(t, "FetchLastRelease")
	p.AssertNotCalled(t, "CacheRelease")
	p.AssertCalled(t, "RestoreCacheRelease")
	assert.Equal(t, "decompression error", err.Error())
}

func TestUpdateChecksumFail(t *testing.T) {
	m := new(mockHttpClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProvider)
	p.On("FetchLastRelease", m).Return(
		&Release{
			Name: "v0.1.2",
		}, nil,
	)
	p.On("CacheRelease", Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("no file error"))
	downloadRelease = func(hcp HttpClientPlugin, r *Release, s string) (string, string, error) {
		return "", "", nil
	}
	funcDecompress = func(src string) (int, error) { return 1, nil }
	checksumRelease = func(binPath, checksumsPath string) error { return fmt.Errorf("checksum error") }
	err := Update(m, p, I18nConf{Verbose: true, Locale: EN}, "14-bis", "0.1.1")

	p.AssertNotCalled(t, "FetchLastRelease")
	p.AssertNotCalled(t, "CacheRelease")
	p.AssertCalled(t, "RestoreCacheRelease")
	assert.Equal(t, "checksum error", err.Error())
}

func TestUpdateInstallationFail(t *testing.T) {
	m := new(mockHttpClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProvider)
	p.On("FetchLastRelease", m).Return(
		&Release{
			Name: "v0.1.2",
		}, nil,
	)
	p.On("CacheRelease", Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("no file error"))
	downloadRelease = func(hcp HttpClientPlugin, r *Release, s string) (string, string, error) {
		return "", "", nil
	}
	funcDecompress = func(src string) (int, error) { return 1, nil }
	checksumRelease = func(binPath, checksumsPath string) error { return nil }
	installRelease = func(srcDir string) error { return fmt.Errorf("installation error") }
	err := Update(m, p, I18nConf{Verbose: true, Locale: EN}, "14-bis", "0.1.1")

	p.AssertNotCalled(t, "FetchLastRelease")
	p.AssertNotCalled(t, "CacheRelease")
	p.AssertCalled(t, "RestoreCacheRelease")
	assert.Equal(t, "installation error", err.Error())
}

func TestUpdate(t *testing.T) {
	m := new(mockHttpClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProvider)
	p.On("FetchLastRelease", m).Return(
		&Release{
			Name: "v0.1.2",
		}, nil,
	)
	p.On("CacheRelease", Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("no file error"))
	downloadRelease = func(hcp HttpClientPlugin, r *Release, s string) (string, string, error) {
		return "", "", nil
	}
	funcDecompress = func(src string) (int, error) { return 1, nil }
	checksumRelease = func(binPath, checksumsPath string) error { return nil }
	installRelease = func(srcDir string) error { return nil }
	err := Update(m, p, I18nConf{Verbose: true, Locale: EN}, "14-bis", "0.1.1")

	p.AssertNotCalled(t, "FetchLastRelease")
	p.AssertNotCalled(t, "CacheRelease")
	p.AssertCalled(t, "RestoreCacheRelease")
	assert.Nil(t, err)
}
