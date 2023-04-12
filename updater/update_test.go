package updater

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

type mockHTTPClientUpdate struct{ mock.Mock }

type mockProviderUpdate struct{ mock.Mock }

func (client *mockHTTPClientUpdate) Do(req *http.Request) (*http.Response, error) {
	args := client.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (provider *mockProviderUpdate) FetchReleases(client pvdr.HTTPClientPlugin) ([]*pvdr.Release, error) {
	args := provider.Called(client)
	return args.Get(0).([]*pvdr.Release), args.Error(1)
}

func (provider *mockProviderUpdate) FetchLastRelease(client pvdr.HTTPClientPlugin) (*pvdr.Release, error) {
	args := provider.Called(client)
	var rel *pvdr.Release
	if args.Get(0) != nil {
		rel, _ = args.Get(0).(*pvdr.Release)
	}

	return rel, args.Error(1)
}

func (provider *mockProviderUpdate) CacheRelease(rel pvdr.Release) error {
	args := provider.Called(rel)
	return args.Error(0)
}

func (provider *mockProviderUpdate) RestoreCacheRelease() (*pvdr.Release, error) {
	args := provider.Called()
	var rel *pvdr.Release
	if args.Get(0) != nil {
		rel, _ = args.Get(0).(*pvdr.Release)
	}

	return rel, args.Error(1)
}

func TestUpdateCheckVersionFail(t *testing.T) {
	m := new(mockHTTPClientUpdate)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(``))),
		}, nil)
	p := new(mockProviderUpdate)
	p.On("FetchLastRelease", m).Return(
		nil, fmt.Errorf("any error"),
	)
	p.On("CacheRelease", pvdr.Release{}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("any error"))

	_, err := UpdateRelease(m, p, "0.0.1", false)
	actual := err.Error()
	expected := "any error"

	p.AssertNotCalled(t, "FetchLastRelease")
	p.AssertNotCalled(t, "CacheRelease")
	p.AssertCalled(t, "RestoreCacheRelease")
	assert.Equal(t, expected, actual)
}

func TestUpdateAlreadyUpToDate(t *testing.T) {
	m := new(mockHTTPClientUpdate)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProviderUpdate)
	p.On("FetchLastRelease", m).Return(
		&pvdr.Release{
			Name: "v0.1.2",
		}, nil,
	)
	p.On("CacheRelease", pvdr.Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("no file error"))

	_, err := UpdateRelease(m, p, "0.1.2", false)
	actual := err.Error()
	expected := "already on the edge"

	p.AssertNotCalled(t, "FetchLastRelease")
	p.AssertNotCalled(t, "CacheRelease")
	p.AssertCalled(t, "RestoreCacheRelease")
	assert.Equal(t, expected, actual)
}

func TestUpdateProcFilePathFail(t *testing.T) {
	m := new(mockHTTPClientUpdate)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProviderUpdate)
	p.On("FetchLastRelease", m).Return(
		&pvdr.Release{
			Name: "v0.1.2",
		}, nil,
	)
	p.On("CacheRelease", pvdr.Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("no file error"))
	mpProcessFilePath = func() (string, error) { return "", fmt.Errorf("process path error") }

	_, err := UpdateRelease(m, p, "0.1.1", false)

	p.AssertNotCalled(t, "FetchLastRelease")
	p.AssertNotCalled(t, "CacheRelease")
	p.AssertCalled(t, "RestoreCacheRelease")
	assert.Equal(t, "process path error", err.Error())
}

func TestUpdateDownloadFail(t *testing.T) {
	m := new(mockHTTPClientUpdate)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProviderUpdate)
	p.On("FetchLastRelease", m).Return(
		&pvdr.Release{
			Name: "v0.1.2",
		}, nil,
	)
	p.On("CacheRelease", pvdr.Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("no file error"))
	mpProcessFilePath = func() (string, error) { return "", nil }
	mpDownloadTo = func(hcp pvdr.HTTPClientPlugin, r *pvdr.Release, s string) (string, string, error) {
		return "", "", fmt.Errorf("download release error")
	}

	_, err := UpdateRelease(m, p, "0.1.1", false)

	p.AssertNotCalled(t, "FetchLastRelease")
	p.AssertNotCalled(t, "CacheRelease")
	p.AssertCalled(t, "RestoreCacheRelease")
	assert.Equal(t, "download release error", err.Error())
}

func TestUpdateDecompressionFail(t *testing.T) {
	m := new(mockHTTPClientUpdate)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProviderUpdate)
	p.On("FetchLastRelease", m).Return(
		&pvdr.Release{
			Name: "v0.1.2",
		}, nil,
	)
	p.On("CacheRelease", pvdr.Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("no file error"))
	mpProcessFilePath = func() (string, error) { return "", nil }
	mpDownloadTo = func(hcp pvdr.HTTPClientPlugin, r *pvdr.Release, s string) (string, string, error) {
		return "", "", nil
	}
	mpDecompress = func(src string) (int, error) { return 0, fmt.Errorf("decompression error") }
	_, err := UpdateRelease(m, p, "0.1.1", false)

	p.AssertNotCalled(t, "FetchLastRelease")
	p.AssertNotCalled(t, "CacheRelease")
	p.AssertCalled(t, "RestoreCacheRelease")
	assert.Equal(t, "decompression error", err.Error())
}

func TestUpdateChecksumFail(t *testing.T) {
	m := new(mockHTTPClientUpdate)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProviderUpdate)
	p.On("FetchLastRelease", m).Return(
		&pvdr.Release{
			Name: "v0.1.2",
		}, nil,
	)
	p.On("CacheRelease", pvdr.Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("no file error"))
	mpProcessFilePath = func() (string, error) { return "", nil }
	mpDownloadTo = func(hcp pvdr.HTTPClientPlugin, r *pvdr.Release, s string) (string, string, error) {
		return "", "", nil
	}
	mpDecompress = func(src string) (int, error) { return 1, nil }
	mpChecksum = func(binPath, checksumsPath string) error { return fmt.Errorf("checksum error") }
	_, err := UpdateRelease(m, p, "0.1.1", false)

	p.AssertNotCalled(t, "FetchLastRelease")
	p.AssertNotCalled(t, "CacheRelease")
	p.AssertCalled(t, "RestoreCacheRelease")
	assert.Equal(t, "checksum error", err.Error())
}

func TestUpdateInstallationFail(t *testing.T) {
	m := new(mockHTTPClientUpdate)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProviderUpdate)
	p.On("FetchLastRelease", m).Return(
		&pvdr.Release{
			Name: "v0.1.2",
		}, nil,
	)
	p.On("CacheRelease", pvdr.Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("no file error"))
	mpProcessFilePath = func() (string, error) { return "", nil }
	mpDownloadTo = func(hcp pvdr.HTTPClientPlugin, r *pvdr.Release, s string) (string, string, error) {
		return "", "", nil
	}
	mpDecompress = func(src string) (int, error) { return 1, nil }
	mpChecksum = func(binPath, checksumsPath string) error { return nil }
	mpInstall = func(srcDir, destDir string) error { return fmt.Errorf("installation error") }
	_, err := UpdateRelease(m, p, "0.1.1", false)

	p.AssertNotCalled(t, "FetchLastRelease")
	p.AssertNotCalled(t, "CacheRelease")
	p.AssertCalled(t, "RestoreCacheRelease")
	assert.Equal(t, "installation error", err.Error())
}

func TestUpdate(t *testing.T) {
	m := new(mockHTTPClientUpdate)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProviderUpdate)
	p.On("FetchLastRelease", m).Return(
		&pvdr.Release{
			Name: "v0.1.2",
		}, nil,
	)
	p.On("CacheRelease", pvdr.Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("no file error"))
	mpProcessFilePath = func() (string, error) { return "/tmp/test-update", nil }
	mpDownloadTo = func(hcp pvdr.HTTPClientPlugin, r *pvdr.Release, s string) (string, string, error) {
		return "", "", nil
	}
	mpDecompress = func(src string) (int, error) { return 1, nil }
	mpChecksum = func(binPath, checksumsPath string) error { return nil }
	mpInstall = func(srcDir, destDir string) error { return nil }
	_, err := UpdateRelease(m, p, "0.1.1", false)

	p.AssertNotCalled(t, "FetchLastRelease")
	p.AssertNotCalled(t, "CacheRelease")
	p.AssertCalled(t, "RestoreCacheRelease")
	assert.Nil(t, err)
}
