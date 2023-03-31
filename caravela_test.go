package caravela

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	httpc "github.com/aureliano/caravela/http"
	"github.com/aureliano/caravela/release"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockHttpClient struct{ mock.Mock }

type mockProvider struct{ mock.Mock }

func (client *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := client.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (provider *mockProvider) FetchReleases(client httpc.HttpClientPlugin) ([]*release.Release, error) {
	args := provider.Called(client)
	return args.Get(0).([]*release.Release), args.Error(1)
}

func (provider *mockProvider) FetchLastRelease(client httpc.HttpClientPlugin) (*release.Release, error) {
	args := provider.Called(client)
	var rel *release.Release
	if args.Get(0) != nil {
		rel = args.Get(0).(*release.Release)
	}

	return rel, args.Error(1)
}

func (provider *mockProvider) CacheRelease(rel release.Release) error {
	args := provider.Called(rel)
	return args.Error(0)
}

func (provider *mockProvider) RestoreCacheRelease() (*release.Release, error) {
	args := provider.Called()
	var rel *release.Release
	if args.Get(0) != nil {
		rel = args.Get(0).(*release.Release)
	}

	return rel, args.Error(1)
}

func TestCheckForUpdatesRestoreCacheCurrentVersionIsOlder(t *testing.T) {
	m := new(mockHttpClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProvider)
	p.On("RestoreCacheRelease").Return(&release.Release{Name: "v0.1.0"}, nil)

	r, _ := CheckForUpdates(m, p, "v0.1.0-alpha")
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
	p.On("RestoreCacheRelease").Return(&release.Release{Name: "v0.1.0"}, nil)

	r, _ := CheckForUpdates(m, p, "v0.1.0")
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
	p.On("CacheRelease", release.Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("any error"))
	p.On("FetchLastRelease", m).Return(
		nil, fmt.Errorf("some error"),
	)

	r, e := CheckForUpdates(m, p, "v0.1.2")
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
	p.On("CacheRelease", release.Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("any error"))
	p.On("FetchLastRelease", m).Return(
		&release.Release{
			Name: "v0.1.2",
		}, nil,
	)

	r, _ := CheckForUpdates(m, p, "v0.1.1")
	assert.Equal(t, r.Name, "v0.1.2")
	p.AssertCalled(t, "FetchLastRelease", m)
	p.AssertCalled(t, "CacheRelease", release.Release{Name: "v0.1.2"})
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
	p.On("CacheRelease", release.Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("any error"))
	p.On("FetchLastRelease", m).Return(
		&release.Release{
			Name: "v0.1.2",
		}, nil,
	)

	r, _ := CheckForUpdates(m, p, "v0.1.2")
	assert.Nil(t, r)
	p.AssertCalled(t, "FetchLastRelease", m)
	p.AssertCalled(t, "CacheRelease", release.Release{Name: "v0.1.2"})
	p.AssertCalled(t, "RestoreCacheRelease")
}

/*
func TestCheckForUpdatesNoCache(t *testing.T) {
	m := new(mockHttpClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProvider)
	p.On("FetchReleases", m).Return(
		[]*release.Release{}, nil,
	)
	p.On("FetchLastRelease", m).Return(
		&release.Release{
			Name: "v0.1.2",
		}, nil,
	)
	p.On("CacheRelease", release.Release{Name: "v0.1.2"}).Return(nil)
	p.On("RestoreCacheRelease").Return(nil, fmt.Errorf("no file error"))

	type testCase struct {
		name     string
		input    string
		expected *release.Release
	}
	testCases := []testCase{
		{
			name:     "current version is the newest",
			input:    "v0.1.2",
			expected: nil,
		}, {
			name:     "current version is NOT the newest",
			input:    "v0.1.1",
			expected: &release.Release{Name: "v0.1.2"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, _ := CheckForUpdates(m, p, tc.input)
			p.AssertCalled(t, "FetchLastRelease", m)
			p.AssertCalled(t, "CacheRelease", release.Release{Name: "v0.1.2"})
			p.AssertCalled(t, "RestoreCacheRelease")

			if (actual == nil && tc.expected != nil) || (actual != nil && tc.expected == nil) ||
				(actual != nil && tc.expected != nil && actual.Name != tc.expected.Name) {
				assert.Fail(t, "expected %v, but got %v", tc.expected, actual)
			}
		})
	}

	rel, err := CheckForUpdates(m, p, "v0.1.1")
	p.AssertCalled(t, "FetchLastRelease", m)
	p.AssertCalled(t, "CacheRelease", release.Release{Name: "v0.1.2"})
	p.AssertCalled(t, "RestoreCacheRelease")
	assert.Nil(t, err, err)

	actual := rel.Name
	expected := "v0.1.2"
	assert.Equal(t, expected, actual)
}

func TestCheckForUpdatesWithCache(t *testing.T) {
	m := new(mockHttpClient)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)
	p := new(mockProvider)
	p.On("FetchReleases", m).Return(
		[]*release.Release{}, nil,
	)
	p.On("FetchLastRelease", m).Return(
		&release.Release{
			Name: "v0.1.2",
		}, nil,
	)
	p.On("CacheRelease", release.Release{}).Return(nil)
	p.On("RestoreCacheRelease").Return(&release.Release{Name: "v0.1.2"}, nil)

	type testCase struct {
		name     string
		input    string
		expected *release.Release
	}
	testCases := []testCase{
		{
			name:     "current version is the newest",
			input:    "v0.1.2",
			expected: nil,
		}, {
			name:     "current version is NOT the newest",
			input:    "v0.1.1",
			expected: &release.Release{Name: "v0.1.2"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, _ := CheckForUpdates(m, p, tc.input)
			p.AssertNotCalled(t, "FetchLastRelease")
			p.AssertNotCalled(t, "CacheRelease")
			p.AssertCalled(t, "RestoreCacheRelease")

			if (actual == nil && tc.expected != nil) || (actual != nil && tc.expected == nil) ||
				(actual != nil && tc.expected != nil && actual.Name != tc.expected.Name) {
				assert.Fail(t, "expected %v, but got %v", tc.expected, actual)
			}
		})
	}

	rel, err := CheckForUpdates(m, p, "v0.1.1")
	p.AssertNotCalled(t, "FetchLastRelease")
	p.AssertNotCalled(t, "CacheRelease")
	p.AssertCalled(t, "RestoreCacheRelease")
	assert.Nil(t, err, err)

	actual := rel.Name
	expected := "v0.1.2"
	assert.Equal(t, expected, actual)
}
*/
