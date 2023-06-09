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

type mockHTTPClientFindUpdate struct{ mock.Mock }

type mockProviderFindUpdate struct{ mock.Mock }

func (client *mockHTTPClientFindUpdate) Do(req *http.Request) (*http.Response, error) {
	args := client.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (provider *mockProviderFindUpdate) FetchReleases(client pvdr.HTTPClientPlugin) ([]*pvdr.Release, error) {
	args := provider.Called(client)
	return args.Get(0).([]*pvdr.Release), args.Error(1)
}

func (provider *mockProviderFindUpdate) FetchLastRelease(client pvdr.HTTPClientPlugin) (*pvdr.Release, error) {
	args := provider.Called(client)
	var rel *pvdr.Release
	if args.Get(0) != nil {
		rel, _ = args.Get(0).(*pvdr.Release)
	}

	return rel, args.Error(1)
}

func (provider *mockProviderFindUpdate) CacheRelease(rel pvdr.Release) error {
	args := provider.Called(rel)
	return args.Error(0)
}

func (provider *mockProviderFindUpdate) RestoreCacheRelease() (*pvdr.Release, error) {
	args := provider.Called()
	var rel *pvdr.Release
	if args.Get(0) != nil {
		rel, _ = args.Get(0).(*pvdr.Release)
	}

	return rel, args.Error(1)
}

func TestCheckUpdatesRestoreCacheCurrentVersionIsOlder(t *testing.T) {
	m := new(mockHTTPClientFindUpdate)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProviderFindUpdate)
	p.On("RestoreCacheRelease").Return(&pvdr.Release{Name: "v0.1.0"}, nil)

	r, _ := FindUpdate(m, p, "v0.1.0-alpha", false)
	assert.Equal(t, r.Name, "v0.1.0")
	p.AssertCalled(t, "RestoreCacheRelease")
}

func TestCheckUpdatesRestoreCacheCurrentVersionOnTheEdge(t *testing.T) {
	m := new(mockHTTPClientFindUpdate)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProviderFindUpdate)
	p.On("RestoreCacheRelease").Return(&pvdr.Release{Name: "v0.1.0"}, nil)

	r, _ := FindUpdate(m, p, "v0.1.0", false)
	assert.Empty(t, r.Name)
	p.AssertCalled(t, "RestoreCacheRelease")
}

func TestCheckUpdatesNoCacheFetchLastReleaseError(t *testing.T) {
	m := new(mockHTTPClientFindUpdate)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProviderFindUpdate)
	p.On("CacheRelease", pvdr.Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("any error"))
	p.On("FetchLastRelease", m).Return(
		nil, fmt.Errorf("some error"),
	)

	r, e := FindUpdate(m, p, "v0.1.2", false)
	assert.Nil(t, r)
	assert.Equal(t, "some error", e.Error())
	p.AssertCalled(t, "FetchLastRelease", m)
	p.AssertCalled(t, "RestoreCacheRelease")
}

func TestCheckUpdatesNoCacheCurrentVersionIsOlder(t *testing.T) {
	m := new(mockHTTPClientFindUpdate)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProviderFindUpdate)
	p.On("CacheRelease", pvdr.Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("any error"))
	p.On("FetchLastRelease", m).Return(
		&pvdr.Release{
			Name: "v0.1.2",
		}, nil,
	)

	r, _ := FindUpdate(m, p, "v0.1.1", false)
	assert.Equal(t, r.Name, "v0.1.2")
	p.AssertCalled(t, "FetchLastRelease", m)
	p.AssertCalled(t, "CacheRelease", pvdr.Release{Name: "v0.1.2"})
	p.AssertCalled(t, "RestoreCacheRelease")
}

func TestCheckUpdatesNoCacheCurrentVersionOnTheEdge(t *testing.T) {
	m := new(mockHTTPClientFindUpdate)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProviderFindUpdate)
	p.On("CacheRelease", pvdr.Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("any error"))
	p.On("FetchLastRelease", m).Return(
		&pvdr.Release{
			Name: "v0.1.2",
		}, nil,
	)

	r, _ := FindUpdate(m, p, "v0.1.2", false)
	assert.Empty(t, r.Name)
	p.AssertCalled(t, "FetchLastRelease", m)
	p.AssertCalled(t, "CacheRelease", pvdr.Release{Name: "v0.1.2"})
	p.AssertCalled(t, "RestoreCacheRelease")
}

func TestCheckUpdatesIgnoreCache(t *testing.T) {
	m := new(mockHTTPClientFindUpdate)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProviderFindUpdate)
	p.On("FetchLastRelease", m).Return(
		&pvdr.Release{
			Name: "v0.1.3",
		}, nil,
	)

	r, _ := FindUpdate(m, p, "v0.1.2", true)
	assert.Equal(t, r.Name, "v0.1.3")
	p.AssertCalled(t, "FetchLastRelease", m)
}

func TestCheckUpdatesIgnoreCacheError(t *testing.T) {
	m := new(mockHTTPClientFindUpdate)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProviderFindUpdate)
	p.On("FetchLastRelease", m).Return(
		nil, fmt.Errorf("some error"),
	)

	r, e := FindUpdate(m, p, "v0.1.2", true)
	assert.Nil(t, r)
	assert.Equal(t, "some error", e.Error())
	p.AssertCalled(t, "FetchLastRelease", m)
}
