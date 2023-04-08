package provider

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSerializeRelease(t *testing.T) {
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

	err := serializeRelease(release)
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

func TestDeserializeRelease(t *testing.T) {
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

	now := time.Now().UTC()
	fname := fmt.Sprintf("release_%s.json", now.Format("2006-01-02"))
	path := filepath.Join(os.TempDir(), fname)

	file, err := os.Create(path)
	assert.Nil(t, err, err)
	_, _ = io.WriteString(file, "{\"name\":\"v0.1.0-dev\",\"description\":\"Development version."+
		"\",\"releasedAt\":\"2023-03-06T09:59:26Z\",\"assets\":[{\"name\":\"f1\",\"url\":\"u1\"},"+
		"{\"name\":\"f2\",\"url\":\"u2\"},{\"name\":\"f3\",\"url\":\"u3\"}]}")
	file.Close()

	actual, err := deserializeRelease()
	assert.Nil(t, err, err)

	assert.Equal(t, release, actual)
	os.Remove(path)
}
