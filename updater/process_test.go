package updater

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessFilePathOsExecFail(t *testing.T) {
	_mpOsExecutable = func() (string, error) { return "", fmt.Errorf("some error") }
	_, err := processFilePath()

	actual := err.Error()
	expected := "error getting running process: some error"

	assert.Equal(t, expected, actual)
}

func TestProcessFilePathUnknownPath(t *testing.T) {
	_mpOsExecutable = func() (string, error) { return "/unknown/path", nil }
	_, err := processFilePath()

	actual := err.Error()
	expected := "error getting information from process: lstat /unknown/path: no such file or directory"

	assert.Equal(t, expected, actual)
}

func TestProcessFilePath(t *testing.T) {
	_mpOsExecutable = os.Executable
	dir, err := processFilePath()
	if err != nil {
		t.Fatal(err)
	}

	parts := strings.Split(dir, string(filepath.Separator))
	if len(parts) < 3 {
		t.Fatal("Expected at least 2 tokens.")
	}

	actual := parts[1]
	expected := filepath.Base(os.TempDir())
	assert.Equal(t, expected, actual)

	actual = parts[2]
	expected = "go-build"
	assert.Contains(t, actual, expected)
}
