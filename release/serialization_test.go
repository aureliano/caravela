package release

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
			Name string
			URL  string
		}{
			{Name: "f1", URL: "u1"},
			{Name: "f2", URL: "u2"},
			{Name: "f3", URL: "u3"},
		},
	}

	err := SerializeRelease(release)
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

func TestDeserializeRelease(t *testing.T) {
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

	now := time.Now().UTC()
	fname := fmt.Sprintf("release_%s.json", now.Format("2006-01-02"))
	path := filepath.Join(os.TempDir(), fname)

	file, err := os.Create(path)
	assert.Nil(t, err, err)
	_, _ = io.WriteString(file, "{\"Name\":\"v0.1.0-dev\",\"Description\":\"Development version.\",\"ReleasedAt\":\"2023-03-06T09:59:26Z\",\"assets\":[{\"Name\":\"f1\",\"URL\":\"u1\"},{\"Name\":\"f2\",\"URL\":\"u2\"},{\"Name\":\"f3\",\"URL\":\"u3\"}]}")
	file.Close()

	actual, err := DeserializeRelease()
	assert.Nil(t, err, err)

	assert.Equal(t, release, actual)
	os.Remove(path)
}
