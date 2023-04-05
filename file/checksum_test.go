package file

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChecksumBinNotFound(t *testing.T) {
	file, err := os.Create(filepath.Join(os.TempDir(), "checksums.txt"))
	assert.Nil(t, err, err)

	_, _ = file.WriteString("dc173fa63edc62745edaa05422a3f2d7413d36b52f10d9e6623a8a946b8792db 14-bis_Linux_x86_64.zip")
	file.Close()

	err = Checksum("/no/file", file.Name())
	assert.NotNil(t, err)
}

func TestChecksumDoesntMatch(t *testing.T) {
	zip := filepath.Join(os.TempDir(), "14-bis_Linux_x86_64")
	err := createZipFile(zip)
	zip = fmt.Sprintf("%s.zip", zip)
	assert.Nil(t, err, err)

	file, err := os.Create(filepath.Join(os.TempDir(), "checksums.txt"))
	assert.Nil(t, err, err)

	_, _ = file.WriteString("12345 14-bis_Linux_x86_64.zip")
	file.Close()

	err = Checksum(zip, file.Name())
	actual := err.Error()
	expected := "checksum failed"
	assert.Equal(t, expected, actual)
}

func TestChecksum(t *testing.T) {
	zip := filepath.Join(os.TempDir(), "14-bis_Linux_x86_64")
	err := createZipFile(zip)
	zip = fmt.Sprintf("%s.zip", zip)
	assert.Nil(t, err, err)

	file, err := os.Create(filepath.Join(os.TempDir(), "checksums.txt"))
	assert.Nil(t, err, err)

	_, _ = file.WriteString("dc173fa63edc62745edaa05422a3f2d7413d36b52f10d9e6623a8a946b8792db 14-bis_Linux_x86_64.zip")
	file.Close()

	err = Checksum(zip, file.Name())
	assert.Nil(t, err, err)
}

func TestGetChecksumFileNotFound(t *testing.T) {
	_, err := getChecksum("/some/path/unexisting.zip", "/some/path/unexisting.txt")
	assert.NotNil(t, err)
}

func TestGetChecksumDoesntMatch(t *testing.T) {
	file, err := os.Create(filepath.Join(os.TempDir(), "checksums.txt"))
	assert.Nil(t, err, err)

	file.WriteString("12345 14-bis_Windows_x86_64.zip")
	file.Close()

	actual, err := getChecksum("14-bis_Linux_x86_64.zip", file.Name())
	expected := ""
	assert.Nil(t, err, err)
	assert.Equal(t, expected, actual)
}

func TestGetChecksum(t *testing.T) {
	file, err := os.Create(filepath.Join(os.TempDir(), "checksums.txt"))
	assert.Nil(t, err, err)

	file.WriteString("12345 14-bis_Linux_x86_64.zip")
	file.Close()

	actual, err := getChecksum("14-bis_Linux_x86_64.zip", file.Name())
	expected := "12345"
	assert.Nil(t, err, err)
	assert.Equal(t, expected, actual)
}
