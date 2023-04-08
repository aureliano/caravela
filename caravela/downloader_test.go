package caravela

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/aureliano/caravela/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type mockHttpPlugin struct{ mock.Mock }

func (plugin *mockHttpPlugin) Do(req *http.Request) (*http.Response, error) {
	args := plugin.Called(req)
	var res *http.Response
	if args.Get(0) != nil {
		res = args.Get(0).(*http.Response)
	}

	return res, args.Error(1)
}

func TestDownloadToBinNotFound(t *testing.T) {
	m := new(mockHttpPlugin)
	m.On("Do", mock.Anything).Return(nil, nil)

	release := new(provider.Release)
	release.Assets = []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}{
		{Name: "14-bis_Xpto_x86_64.tar.gz", URL: "http://file-linux.tar.gz"},
		{Name: "14-bis_Lorem_x86_64.zip", URL: "http://file-windows.zip"},
		{Name: "14-bis_Ipsum_x86_64.tar.gz", URL: "http://file-darwin.tar.gz"},
		{Name: "checksums.txt", URL: "http://checksums.txt"},
	}

	_, _, err := downloadTo(m, release, "")
	assert.Contains(t, err.Error(), "there is no version compatible with")
	m.AssertNotCalled(t, "Do", mock.Anything)
}

func TestDownloadToChecksumsNotFound(t *testing.T) {
	m := new(mockHttpPlugin)
	m.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte("12345"))),
	}, nil)

	release := new(provider.Release)
	release.Assets = []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}{
		{Name: "14-bis_Linux_x86_64.tar.gz", URL: "http://file-linux.tar.gz"},
		{Name: "14-bis_Windows_x86_64.zip", URL: "http://file-windows.zip"},
		{Name: "14-bis_Darwin_x86_64.tar.gz", URL: "http://file-darwin.tar.gz"},
	}

	_, _, err := downloadTo(m, release, os.TempDir())
	assert.Contains(t, err.Error(), "file checksums.txt not found")
	m.AssertCalled(t, "Do", mock.Anything)
}

func TestDownloadDownloadBinError(t *testing.T) {
	m := new(mockHttpPlugin)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("12345"))),
		}, nil)

	release := new(provider.Release)
	release.Assets = []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}{
		{Name: "14-bis_Linux_x86_64.tar.gz", URL: "http://file-linux.tar.gz"},
		{Name: "14-bis_Windows_x86_64.zip", URL: "http://file-windows.zip"},
		{Name: "14-bis_Darwin_x86_64.tar.gz", URL: "http://file-darwin.tar.gz"},
		{Name: "checksums.txt", URL: "http://checksums.txt"},
	}

	mpDownloadFile = func(client provider.HTTPClientPlugin, sourceUrl, dest string) error {
		if strings.Contains(dest, "14-bis_") {
			return fmt.Errorf("failed to download binary")
		}

		return nil
	}

	_, _, err := downloadTo(m, release, os.TempDir())
	assert.Equal(t, "failed to download binary", err.Error())
}

func TestDownloadDownloadChecksumsError(t *testing.T) {
	m := new(mockHttpPlugin)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("12345"))),
		}, nil)

	release := new(provider.Release)
	release.Assets = []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}{
		{Name: "14-bis_Linux_x86_64.tar.gz", URL: "http://file-linux.tar.gz"},
		{Name: "14-bis_Windows_x86_64.zip", URL: "http://file-windows.zip"},
		{Name: "14-bis_Darwin_x86_64.tar.gz", URL: "http://file-darwin.tar.gz"},
		{Name: "checksums.txt", URL: "http://checksums.txt"},
	}

	mpDownloadFile = func(client provider.HTTPClientPlugin, sourceUrl, dest string) error {
		if strings.Contains(dest, "checksums.txt") {
			return fmt.Errorf("failed to download checksums")
		}

		return nil
	}

	_, _, err := downloadTo(m, release, os.TempDir())
	assert.Equal(t, "failed to download checksums", err.Error())
}

func TestDownloadRelease(t *testing.T) {
	m := new(mockHttpPlugin)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("12345"))),
		}, nil)

	release := new(provider.Release)
	release.Assets = []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}{
		{Name: "14-bis_Linux_x86_64.tar.gz", URL: "http://file-linux.tar.gz"},
		{Name: "14-bis_Windows_x86_64.zip", URL: "http://file-windows.zip"},
		{Name: "14-bis_Darwin_x86_64.tar.gz", URL: "http://file-darwin.tar.gz"},
		{Name: "checksums.txt", URL: "http://checksums.txt"},
	}

	osName := cases.Title(language.Und).String(runtime.GOOS)
	var suffix string
	if osName == "Windows" {
		suffix = "zip"
	} else {
		suffix = "tar.gz"
	}

	mpDownloadFile = downloadFile

	dir := os.TempDir()
	abin, achecksum, err := downloadTo(m, release, dir)
	ebin, echecksum := filepath.Join(dir, fmt.Sprintf("14-bis_%s_x86_64.%s", osName, suffix)), filepath.Join(dir, "checksums.txt")

	assert.Nil(t, err, err)
	assert.Equal(t, ebin, abin)
	assert.Equal(t, echecksum, achecksum)
	m.AssertCalled(t, "Do", mock.Anything)
}

func TestDownloadFileWrongDest(t *testing.T) {
	m := new(mockHttpPlugin)
	m.On("Do", mock.Anything).Return(nil, nil)

	err := downloadFile(m, "http://file-linux.tar.gz", filepath.Join("unknown", "path"))
	assert.NotNil(t, err, err)
	m.AssertNotCalled(t, "Do", mock.Anything)
}

func TestDownloadFileHttpError(t *testing.T) {
	m := new(mockHttpPlugin)
	m.On("Do", mock.Anything).Return(nil, fmt.Errorf("some error"))

	dest := filepath.Join(os.TempDir(), "file-linux.tar.gz")
	actual := downloadFile(m, "http://file-linux.tar.gz", dest)
	expected := "some error"

	assert.Equal(t, expected, actual.Error())
	m.AssertCalled(t, "Do", mock.Anything)
}

func TestDownloadFile404(t *testing.T) {
	m := new(mockHttpPlugin)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewReader([]byte("not found"))),
		}, nil)

	dest := filepath.Join(os.TempDir(), "file-linux.tar.gz")
	actual := downloadFile(m, "http://file-linux.tar.gz", dest)
	expected := fmt.Errorf("http error (404)")

	assert.Equal(t, expected.Error(), actual.Error())
	m.AssertCalled(t, "Do", mock.Anything)
}

func TestDownloadFile(t *testing.T) {
	m := new(mockHttpPlugin)
	m.On("Do", mock.Anything).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("12345"))),
		}, nil)
	release := new(provider.Release)
	release.Assets = []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}{
		{Name: "14-bis_Linux_x86_64.tar.gz", URL: "http://file-linux.tar.gz"},
		{Name: "14-bis_Windows_x86_64.zip", URL: "http://file-windows.zip"},
		{Name: "14-bis_Darwin_x86_64.tar.gz", URL: "http://file-darwin.tar.gz"},
		{Name: "checksums.txt", URL: "http://checksums.txt"},
	}

	dest := filepath.Join(os.TempDir(), "file-linux.tar.gz")
	err := downloadFile(m, "http://file-linux.tar.gz", dest)
	assert.Nil(t, err, err)
	m.AssertCalled(t, "Do", mock.Anything)

	file, err := os.Open(dest)
	assert.Nil(t, err, err)

	data, err := io.ReadAll(file)
	assert.Nil(t, err, err)

	actual := string(data)
	expected := "12345"
	assert.Equal(t, expected, actual)

	os.Remove(file.Name())
}

func TestFetchReleaseFileUrl(t *testing.T) {
	release := new(provider.Release)
	release.Assets = []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}{
		{Name: "14-bis_Linux_x86_64.tar.gz", URL: "http://file-linux.tar.gz"},
		{Name: "14-bis_Windows_x86_64.zip", URL: "http://file-windows.zip"},
		{Name: "14-bis_Darwin_x86_64.tar.gz", URL: "http://file-darwin.tar.gz"},
		{Name: "checksums.txt", URL: "http://checksums.txt"},
	}

	type testCase struct {
		name     string
		input    string
		expected []string
	}
	testCases := []testCase{
		{
			name:     "linux os",
			input:    "linux",
			expected: []string{"14-bis_Linux_x86_64.tar.gz", "http://file-linux.tar.gz"},
		},
		{
			name:     "windows os",
			input:    "windows",
			expected: []string{"14-bis_Windows_x86_64.zip", "http://file-windows.zip"},
		},
		{
			name:     "darwin os",
			input:    "darwin",
			expected: []string{"14-bis_Darwin_x86_64.tar.gz", "http://file-darwin.tar.gz"},
		},
		{
			name:     "unknown os",
			input:    "unknown",
			expected: []string{"", ""},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			name, url := findReleaseFileUrl(tc.input, release)
			if name != tc.expected[0] && url != tc.expected[1] {
				t.Errorf("expected [%s,  %s], but got [%s, %s]", tc.expected[0], tc.expected[1], name, url)
			}
		})
	}
}

func TestFindChecksumsFileUrlNoCheckSums(t *testing.T) {
	release := new(provider.Release)
	release.Assets = []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}{
		{Name: "14-bis_Linux_x86_64.tar.gz", URL: "http://file-linux.tar.gz"},
		{Name: "14-bis_Windows_x86_64.zip", URL: "http://file-windows.zip"},
		{Name: "14-bis_Darwin_x86_64.tar.gz", URL: "http://file-darwin.tar.gz"},
	}

	actual := findChecksumsFileUrl(release)
	expected := ""
	assert.Equal(t, expected, actual)
}

func TestFindChecksumsFileUrl(t *testing.T) {
	release := new(provider.Release)
	release.Assets = []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}{
		{Name: "14-bis_Linux_x86_64.tar.gz", URL: "http://file-linux.tar.gz"},
		{Name: "14-bis_Windows_x86_64.zip", URL: "http://file-windows.zip"},
		{Name: "14-bis_Darwin_x86_64.tar.gz", URL: "http://file-darwin.tar.gz"},
		{Name: "checksums.txt", URL: "http://checksums.txt"},
	}

	actual := findChecksumsFileUrl(release)
	expected := "http://checksums.txt"
	assert.Equal(t, expected, actual)
}
