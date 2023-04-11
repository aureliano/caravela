package provider

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

type mockGithubDecorator struct{ mock.Mock }

func (decorator *mockGithubDecorator) Do(req *http.Request) (*http.Response, error) {
	args := decorator.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestGithubFetchLastReleaseValidationError(t *testing.T) {
	m := new(mockDecorator)
	p := GithubProvider{}

	_, err := p.FetchLastRelease(m)
	assert.Equal(t, "host is required", err.Error())
}

func TestGithubFetchLastReleaseErrorOnFetchReleases(t *testing.T) {
	m := new(mockDecorator)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(bytes.NewReader([]byte(``))),
		}, nil)

	provider := GithubProvider{Host: "github.com", Port: 80, ProjectPath: "massis/oalienista"}
	_, err := provider.FetchLastRelease(m)
	m.AssertCalled(t, "Do", mock.Anything)
	assert.Equal(t, err.Error(), "github integration error: 500")
}

func TestGithubFetchLastRelease(t *testing.T) {
	m := new(mockDecorator)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(bytes.NewReader(
				[]byte(`[{"tag_name":"v0.1.0"},{"tag_name":"v0.1.1"},{"tag_name":"v0.1.2"}]`))),
		}, nil)

	provider := GithubProvider{Host: "github.com", Port: 80, ProjectPath: "massis/oalienista"}
	actual, err := provider.FetchLastRelease(m)
	m.AssertCalled(t, "Do", mock.Anything)
	assert.Nil(t, err, err)
	assert.Equal(t, "v0.1.2", actual.Name)
}

func TestGithubCacheRelease(t *testing.T) {
	release := &Release{
		Name:        "v0.1.0-dev",
		Description: "Development version.",
		ReleasedAt:  time.Date(2023, 3, 6, 9, 59, 26, 0, time.UTC),
		Assets: []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}{
			{Name: "f1", URL: "u1"},
			{Name: "f2", URL: "u2"},
			{Name: "f3", URL: "u3"},
		},
	}
	provider := GithubProvider{}

	err := provider.CacheRelease(*release)
	assert.Nil(t, err, err)

	now := time.Now().UTC()
	fname := fmt.Sprintf("release_%s.json", now.Format("2006-01-02"))

	file := filepath.Join(os.TempDir(), fname)
	bytes, err := os.ReadFile(file)
	assert.Nil(t, err, err)

	json := string(bytes)
	assert.Equal(t, "{\"name\":\"v0.1.0-dev\",\"description\":\"Development version.\","+
		"\"releasedAt\":\"2023-03-06T09:59:26Z\",\"assets\":[{\"name\":\"f1\",\"url\":\"u1\"},"+
		"{\"name\":\"f2\",\"url\":\"u2\"},{\"name\":\"f3\",\"url\":\"u3\"}]}", json)

	os.Remove(file)
}

func TestGithubRestoreCacheRelease(t *testing.T) {
	release := &Release{
		Name:        "v0.1.0-dev",
		Description: "Development version.",
		ReleasedAt:  time.Date(2023, 3, 6, 9, 59, 26, 0, time.UTC),
		Assets: []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}{
			{Name: "f1", URL: "u1"},
			{Name: "f2", URL: "u2"},
			{Name: "f3", URL: "u3"},
		},
	}
	provider := GithubProvider{}

	now := time.Now().UTC()
	fname := fmt.Sprintf("release_%s.json", now.Format("2006-01-02"))
	path := filepath.Join(os.TempDir(), fname)

	file, err := os.Create(path)
	assert.Nil(t, err, err)
	_, _ = io.WriteString(file, "{\"name\":\"v0.1.0-dev\",\"description\":\"Development version.\","+
		"\"releasedAt\":\"2023-03-06T09:59:26Z\",\"Assets\":[{\"name\":\"f1\",\"url\":\"u1\"},"+
		"{\"name\":\"f2\",\"url\":\"u2\"},{\"name\":\"f3\",\"url\":\"u3\"}]}")
	file.Close()

	actual, err := provider.RestoreCacheRelease()
	assert.Nil(t, err, err)

	assert.Equal(t, release, actual)
	os.Remove(path)
}

func TestGithubReleaseCompareTo(t *testing.T) {
	r1 := &GithubRelease{Name: "v0.1.0"}
	r2 := &GithubRelease{Name: "v0.1.1"}

	assert.Equal(t, r1.CompareTo(r2), -1)
}

func TestFetchGithubReleasesEmpty(t *testing.T) {
	m := new(mockDecorator)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
		}, nil)

	provider := GithubProvider{}
	actual, err := fetchGithubReleases(provider, m)
	expected := []*GithubRelease{}

	assert.Nil(t, err, err)
	m.AssertCalled(t, "Do", mock.Anything)

	s1, s2 := len(actual), len(expected)
	if s1 != s2 {
		assert.Fail(t, "expected %d == %d", s1, s2)
	} else if s1 > 0 {
		assert.Fail(t, "expected 0, but got %d", s1)
	}
}

func TestFetchGithubReleasesError(t *testing.T) {
	m := new(mockDecorator)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       nil,
		}, errors.New("http test"))

	provider := GithubProvider{}
	actual, err := fetchGithubReleases(provider, m)

	m.AssertCalled(t, "Do", mock.Anything)
	assert.Equal(t, err.Error(), "http test")
	assert.Nil(t, actual)
}

func TestFetchGithubReleasesBrokenJson(t *testing.T) {
	m := new(mockDecorator)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[{]`))),
		}, nil)

	provider := GithubProvider{}
	actual, err := fetchGithubReleases(provider, m)

	m.AssertCalled(t, "Do", mock.Anything)
	assert.NotNil(t, err)
	assert.Equal(t, len(actual), 0)
}

func TestFetchGithubReleases(t *testing.T) {
	m := new(mockGithubDecorator)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(
				bytes.NewReader([]byte(`[{"tag_name":"v0.1.0"},{"tag_name":"v0.1.1"},{"tag_name":"v0.1.2"}]`))),
		}, nil)

	provider := GithubProvider{}
	actual, err := fetchGithubReleases(provider, m)
	expected := []*GithubRelease{
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

func TestBuildGithubServiceUrl(t *testing.T) {
	p := GithubProvider{"www.domain.com.br", 80, false, "aureliano/caravela", time.Second * 30}
	expected := "http://www.domain.com.br:80/repos/aureliano/caravela/releases"
	actual := buildGithubServiceURL(p)

	assert.Equal(t, expected, actual)
}

func TestBuildGithubServiceUrlSsl(t *testing.T) {
	p := GithubProvider{"www.domain.com.br", 80, true, "aureliano/caravela", time.Second * 30}
	expected := "https://www.domain.com.br:80/repos/aureliano/caravela/releases"
	actual := buildGithubServiceURL(p)

	assert.Equal(t, expected, actual)
}

func TestConvertGithubReleases(t *testing.T) {
	g1 := &GithubRelease{
		Name:        "v0.1.0",
		Body:        "Unit test.",
		PublishedAt: time.Date(2023, 3, 9, 14, 11, 18, 0, time.UTC),
	}

	g1.Assets = []struct {
		Name string "json:\"name\""
		URL  string "json:\"browser_download_url\""
	}{
		{Name: "14-bis_Linux_x86_64.tar.gz", URL: "http://file-linux.tar.gz"},
		{Name: "14-bis_Windows_x86_64.zip", URL: "http://file-windows.zip"},
		{Name: "14-bis_Darwin_x86_64.tar.gz", URL: "http://file-darwin.tar.gz"},
		{Name: "checksums.txt", URL: "http://checksums.txt"},
	}

	g2 := &GithubRelease{
		Name:        "v0.1.1",
		Body:        "Bug fix.",
		PublishedAt: time.Date(2023, 3, 9, 14, 30, 21, 0, time.UTC),
	}

	g2.Assets = []struct {
		Name string "json:\"name\""
		URL  string "json:\"browser_download_url\""
	}{
		{Name: "qtbis_Linux_x86_64.tar.gz", URL: "http://file-lnx.tar.gz"},
		{Name: "qtbis_Windows_x86_64.zip", URL: "http://file-wdws.zip"},
		{Name: "qtbis_Darwin_x86_64.tar.gz", URL: "http://file-drwn.tar.gz"},
		{Name: "checksums.txt", URL: "http://checksums.txt"},
	}

	sources := []*GithubRelease{g1, g2}
	target := convertGithubReleases(sources)

	for i, source := range sources {
		release := target[i]
		assert.Equal(t, source.Name, release.Name)
		assert.Equal(t, source.Body, release.Description)
		assert.Equal(t, source.PublishedAt, release.ReleasedAt)

		for i, expected := range source.Assets {
			actual := release.Assets[i]
			assert.Equal(t, expected.Name, actual.Name)
			assert.Equal(t, expected.URL, actual.URL)
		}
	}
}

func TestConvertGithubToBase(t *testing.T) {
	g := &GithubRelease{
		Name:        "v0.1.0",
		Body:        "Unit test.",
		PublishedAt: time.Date(2023, 3, 9, 14, 11, 18, 0, time.UTC),
	}

	g.Assets = []struct {
		Name string "json:\"name\""
		URL  string "json:\"browser_download_url\""
	}{
		{Name: "14-bis_Linux_x86_64.tar.gz", URL: "http://file-linux.tar.gz"},
		{Name: "14-bis_Windows_x86_64.zip", URL: "http://file-windows.zip"},
		{Name: "14-bis_Darwin_x86_64.tar.gz", URL: "http://file-darwin.tar.gz"},
		{Name: "checksums.txt", URL: "http://checksums.txt"},
	}

	r := convertGithubToBase(g)
	assert.Equal(t, r.Name, g.Name)
	assert.Equal(t, r.Description, g.Body)
	assert.Equal(t, r.ReleasedAt, g.PublishedAt)

	for i, expected := range g.Assets {
		actual := r.Assets[i]
		assert.Equal(t, expected.Name, actual.Name)
		assert.Equal(t, expected.URL, actual.URL)
	}
}

func TestValidateGithubProviderInvalidHost(t *testing.T) {
	p := GithubProvider{Host: "", Port: 80, ProjectPath: "massis/oalienista"}
	expected := "host is required"
	actual := validateGithubProvider(p).Error()

	assert.Equal(t, expected, actual)
}

func TestValidateGithubProviderInvalidPort(t *testing.T) {
	p := GithubProvider{Host: "Github.com", Port: 0, ProjectPath: "massis/oalienista"}
	expected := "port must be > 0"
	actual := validateGithubProvider(p).Error()

	assert.Equal(t, expected, actual)
}

func TestValidateGithubProviderInvalid(t *testing.T) {
	p := GithubProvider{Host: "Github.com", Port: 80, ProjectPath: ""}
	expected := "project path is required"
	actual := validateGithubProvider(p).Error()

	assert.Equal(t, expected, actual)
}

func TestValidateGithubProvider(t *testing.T) {
	p := GithubProvider{Host: "Github.com", Port: 80, ProjectPath: "massis/oalienista"}
	assert.Nil(t, validateGithubProvider(p))
}
func TestInitGithubProviderHttp(t *testing.T) {
	p := GithubProvider{Port: 0, Ssl: false}
	initGithubProvider(&p)

	assert.Equal(t, uint(80), p.Port)
	assert.False(t, p.Ssl)

	p.Ssl = true
	initGithubProvider(&p)
	assert.Equal(t, uint(80), p.Port)
	assert.False(t, p.Ssl)
}

func TestInitGithubProviderHttps(t *testing.T) {
	p := GithubProvider{Port: 0, Ssl: true}
	initGithubProvider(&p)

	assert.Equal(t, uint(443), p.Port)
	assert.True(t, p.Ssl)

	p.Ssl = false
	initGithubProvider(&p)
	assert.Equal(t, uint(443), p.Port)
	assert.True(t, p.Ssl)
}

func TestInitGithubProviderTimeout(t *testing.T) {
	p := GithubProvider{Port: 8443, Ssl: true, Timeout: 0}
	initGithubProvider(&p)

	assert.Equal(t, time.Second*30, p.Timeout)

	p.Timeout = time.Second * 57
	initGithubProvider(&p)

	assert.Equal(t, time.Second*57, p.Timeout)
}
