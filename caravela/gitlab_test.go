package caravela

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDecorator struct{ mock.Mock }

func (decorator *mockDecorator) Do(req *http.Request) (*http.Response, error) {
	args := decorator.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestFetchLastReleaseValidationError(t *testing.T) {
	m := new(mockDecorator)
	p := GitlabProvider{}

	_, err := p.FetchLastRelease(m)
	assert.Equal(t, "host is required", err.Error())
}

func TestFetchLastReleaseErrorOnFetchReleases(t *testing.T) {
	m := new(mockDecorator)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 500,
			Body:       io.NopCloser(bytes.NewReader([]byte(``))),
		}, nil)

	provider := GitlabProvider{Host: "gitlab.com", Port: 80, ProjectPath: "massis/oalienista"}
	_, err := provider.FetchLastRelease(m)
	m.AssertCalled(t, "Do", mock.Anything)
	assert.Equal(t, err.Error(), "Gitlab integration error: 500")
}

func TestFetchLastRelease(t *testing.T) {
	m := new(mockDecorator)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[{"tag_name":"v0.1.0"},{"tag_name":"v0.1.1"},{"tag_name":"v0.1.2"}]`))),
		}, nil)

	provider := GitlabProvider{Host: "gitlab.com", Port: 80, ProjectPath: "massis/oalienista"}
	actual, err := provider.FetchLastRelease(m)
	m.AssertCalled(t, "Do", mock.Anything)
	assert.Nil(t, err, err)
	assert.Equal(t, "v0.1.2", actual.Name)
}

func TestCacheRelease(t *testing.T) {
	release := &Release{
		Name:        "v0.1.0-dev",
		Description: "Development version.",
		ReleasedAt:  time.Date(2023, 3, 6, 9, 59, 26, 0, time.UTC),
		Assets: []struct {
			Name string
			URL  string
		}{
			{Name: "f1", URL: "u1"},
			{Name: "f2", URL: "u2"},
			{Name: "f3", URL: "u3"},
		},
	}
	provider := GitlabProvider{}

	err := provider.CacheRelease(*release)
	assert.Nil(t, err, err)

	now := time.Now().UTC()
	fname := fmt.Sprintf("release_%s.json", now.Format("2006-01-02"))

	file := filepath.Join(os.TempDir(), fname)
	bytes, err := os.ReadFile(file)
	assert.Nil(t, err, err)

	json := string(bytes)
	assert.Equal(t, "{\"Name\":\"v0.1.0-dev\",\"Description\":\"Development version.\",\"ReleasedAt\":\"2023-03-06T09:59:26Z\",\"Assets\":[{\"Name\":\"f1\",\"URL\":\"u1\"},{\"Name\":\"f2\",\"URL\":\"u2\"},{\"Name\":\"f3\",\"URL\":\"u3\"}]}", json)

	os.Remove(file)
}

func TestRestoreCacheRelease(t *testing.T) {
	release := &Release{
		Name:        "v0.1.0-dev",
		Description: "Development version.",
		ReleasedAt:  time.Date(2023, 3, 6, 9, 59, 26, 0, time.UTC),
		Assets: []struct {
			Name string
			URL  string
		}{
			{Name: "f1", URL: "u1"},
			{Name: "f2", URL: "u2"},
			{Name: "f3", URL: "u3"},
		},
	}
	provider := GitlabProvider{}

	now := time.Now().UTC()
	fname := fmt.Sprintf("release_%s.json", now.Format("2006-01-02"))
	path := filepath.Join(os.TempDir(), fname)

	file, err := os.Create(path)
	assert.Nil(t, err, err)
	_, _ = io.WriteString(file, "{\"Name\":\"v0.1.0-dev\",\"Description\":\"Development version.\",\"ReleasedAt\":\"2023-03-06T09:59:26Z\",\"Assets\":[{\"Name\":\"f1\",\"URL\":\"u1\"},{\"Name\":\"f2\",\"URL\":\"u2\"},{\"Name\":\"f3\",\"URL\":\"u3\"}]}")
	file.Close()

	actual, err := provider.RestoreCacheRelease()
	assert.Nil(t, err, err)

	assert.Equal(t, release, actual)
	os.Remove(path)
}

func TestGitlabReleaseCompareTo(t *testing.T) {
	r1 := &GitlabRelease{Name: "v0.1.0"}
	r2 := &GitlabRelease{Name: "v0.1.1"}

	assert.Equal(t, r1.CompareTo(r2), -1)
}

func TestFetchReleasesEmpty(t *testing.T) {
	m := new(mockDecorator)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)

	provider := GitlabProvider{}
	actual, err := fetchReleases(provider, m)
	expected := []*GitlabRelease{}

	assert.Nil(t, err, err)
	m.AssertCalled(t, "Do", mock.Anything)

	s1, s2 := len(actual), len(expected)
	if s1 != s2 {
		assert.Fail(t, "expected %d == %d", s1, s2)
	} else if s1 > 0 {
		assert.Fail(t, "expected 0, but got %d", s1)
	}
}

func TestFetchReleasesError(t *testing.T) {
	m := new(mockDecorator)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 400,
			Body:       nil,
		}, errors.New("http test"))

	provider := GitlabProvider{}
	actual, err := fetchReleases(provider, m)

	m.AssertCalled(t, "Do", mock.Anything)
	assert.Equal(t, err.Error(), "http test")
	assert.Nil(t, actual)
}

func TestFetchReleasesBrokenJson(t *testing.T) {
	m := new(mockDecorator)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[{]`))),
		}, nil)

	provider := GitlabProvider{}
	actual, err := fetchReleases(provider, m)

	m.AssertCalled(t, "Do", mock.Anything)
	assert.NotNil(t, err)
	assert.Equal(t, len(actual), 0)
}

func TestFetchReleases(t *testing.T) {
	m := new(mockDecorator)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[{"tag_name":"v0.1.0"},{"tag_name":"v0.1.1"},{"tag_name":"v0.1.2"}]`))),
		}, nil)

	provider := GitlabProvider{}
	actual, err := fetchReleases(provider, m)
	expected := []*GitlabRelease{
		{Name: "v0.1.0"},
		{Name: "v0.1.1"},
		{Name: "v0.1.2"},
	}

	assert.Nil(t, err, err)
	m.AssertCalled(t, "Do", mock.Anything)

	s1, s2 := len(actual), len(expected)
	assert.Equal(t, s1, s2)

	for i := 0; i < s1; i++ {
		assert.Equal(t, expected[i].Name, actual[i].Name)
	}
}

func TestBuildServiceUrl(t *testing.T) {
	p := GitlabProvider{"www.domain.com.br", 80, false, "aureliano/caravela"}
	expected := "http://www.domain.com.br:80/api/v4/projects/aureliano%2Fcaravela/releases"
	actual := buildServiceUrl(p)

	assert.Equal(t, expected, actual)
}

func TestBuildServiceUrlSsl(t *testing.T) {
	p := GitlabProvider{"www.domain.com.br", 80, true, "aureliano/caravela"}
	expected := "https://www.domain.com.br:80/api/v4/projects/aureliano%2Fcaravela/releases"
	actual := buildServiceUrl(p)

	assert.Equal(t, expected, actual)
}

func TestConvertReleases(t *testing.T) {
	g1 := &GitlabRelease{
		Name:        "v0.1.0",
		Description: "Unit test.",
		ReleaseAt:   time.Date(2023, 3, 9, 14, 11, 18, 0, time.UTC),
	}

	g1.Assets.Links = []struct {
		Name string "json:\"name\""
		URL  string "json:\"url\""
	}{
		{Name: "14-bis_Linux_x86_64.tar.gz", URL: "http://file-linux.tar.gz"},
		{Name: "14-bis_Windows_x86_64.zip", URL: "http://file-windows.zip"},
		{Name: "14-bis_Darwin_x86_64.tar.gz", URL: "http://file-darwin.tar.gz"},
		{Name: "checksums.txt", URL: "http://checksums.txt"},
	}

	g2 := &GitlabRelease{
		Name:        "v0.1.1",
		Description: "Bug fix.",
		ReleaseAt:   time.Date(2023, 3, 9, 14, 30, 21, 0, time.UTC),
	}

	g2.Assets.Links = []struct {
		Name string "json:\"name\""
		URL  string "json:\"url\""
	}{
		{Name: "qtbis_Linux_x86_64.tar.gz", URL: "http://file-lnx.tar.gz"},
		{Name: "qtbis_Windows_x86_64.zip", URL: "http://file-wdws.zip"},
		{Name: "qtbis_Darwin_x86_64.tar.gz", URL: "http://file-drwn.tar.gz"},
		{Name: "checksums.txt", URL: "http://checksums.txt"},
	}

	sources := []*GitlabRelease{g1, g2}
	target := convertReleases(sources)

	for i, source := range sources {
		release := target[i]
		assert.Equal(t, source.Name, release.Name)
		assert.Equal(t, source.Description, release.Description)
		assert.Equal(t, source.ReleaseAt, release.ReleasedAt)

		for i, expected := range source.Assets.Links {
			actual := release.Assets[i]
			assert.Equal(t, expected.Name, actual.Name)
			assert.Equal(t, expected.URL, actual.URL)
		}
	}
}

func TestConvertToBase(t *testing.T) {
	g := &GitlabRelease{
		Name:        "v0.1.0",
		Description: "Unit test.",
		ReleaseAt:   time.Date(2023, 3, 9, 14, 11, 18, 0, time.UTC),
	}

	g.Assets.Links = []struct {
		Name string "json:\"name\""
		URL  string "json:\"url\""
	}{
		{Name: "14-bis_Linux_x86_64.tar.gz", URL: "http://file-linux.tar.gz"},
		{Name: "14-bis_Windows_x86_64.zip", URL: "http://file-windows.zip"},
		{Name: "14-bis_Darwin_x86_64.tar.gz", URL: "http://file-darwin.tar.gz"},
		{Name: "checksums.txt", URL: "http://checksums.txt"},
	}

	r := convertToBase(g)
	assert.Equal(t, r.Name, g.Name)
	assert.Equal(t, r.Description, g.Description)
	assert.Equal(t, r.ReleasedAt, g.ReleaseAt)

	for i, expected := range g.Assets.Links {
		actual := r.Assets[i]
		assert.Equal(t, expected.Name, actual.Name)
		assert.Equal(t, expected.URL, actual.URL)
	}
}

func TestValidateProviderInvalidHost(t *testing.T) {
	p := GitlabProvider{Host: "", Port: 80, ProjectPath: "massis/oalienista"}
	expected := "host is required"
	actual := validateProvider(p).Error()

	assert.Equal(t, expected, actual)
}

func TestValidateProviderInvalidPort(t *testing.T) {
	p := GitlabProvider{Host: "gitlab.com", Port: 0, ProjectPath: "massis/oalienista"}
	expected := "port must be > 0"
	actual := validateProvider(p).Error()

	assert.Equal(t, expected, actual)
}

func TestValidateProviderInvalid(t *testing.T) {
	p := GitlabProvider{Host: "gitlab.com", Port: 80, ProjectPath: ""}
	expected := "project path is required"
	actual := validateProvider(p).Error()

	assert.Equal(t, expected, actual)
}

func TestValidateProvider(t *testing.T) {
	p := GitlabProvider{Host: "gitlab.com", Port: 80, ProjectPath: "massis/oalienista"}
	assert.Nil(t, validateProvider(p))
}
