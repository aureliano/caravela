package caravela

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecompressUnsupportedType(t *testing.T) {
	done, err := decompress("file.rar")
	assert.True(t, done == 0)

	actual := err.Error()
	expected := ".rar not supported for decompression"
	assert.Equal(t, expected, actual)
}

func TestDecompressZip(t *testing.T) {
	zip := filepath.Join(os.TempDir(), "14-bis_Linux_x86_64")
	err := createZipFile(zip)
	zip = fmt.Sprintf("%s.zip", zip)
	assert.Nil(t, err, err)

	actual, err := decompress(zip)
	expected := 2

	assert.Nil(t, err, err)
	assert.Equal(t, expected, actual)

	os.Remove(zip)
}

func TestDecompressUngzip(t *testing.T) {
	tgz := filepath.Join(os.TempDir(), "14-bis_Linux_x86_64")
	err := createGZipFile(tgz)
	tgz = fmt.Sprintf("%s.tar.gz", tgz)
	assert.Nil(t, err, err)

	actual, err := decompress(tgz)
	expected := 2
	assert.Nil(t, err, err)
	assert.Equal(t, expected, actual)

	os.Remove(tgz)
}

func TestUnzipInvalidSrc(t *testing.T) {
	_, err := unzip(filepath.Join("unknown", "path"))
	assert.NotNil(t, err)
}

func TestUnzip(t *testing.T) {
	zip := filepath.Join(os.TempDir(), "14-bis_Linux_x86_64")
	err := createZipFile(zip)
	zip = fmt.Sprintf("%s.zip", zip)
	assert.Nil(t, err, err)

	actual, err := unzip(zip)
	expected := 2

	assert.Nil(t, err, err)
	assert.Equal(t, expected, actual)

	os.Remove(zip)
}

func TestUngzipInvalidSrc(t *testing.T) {
	_, err := ungzip(filepath.Join("unknown", "path"))
	assert.NotNil(t, err)
}

func TestUngzip(t *testing.T) {
	tgz := filepath.Join(os.TempDir(), "14-bis_Linux_x86_64")
	err := createGZipFile(tgz)
	tgz = fmt.Sprintf("%s.tar.gz", tgz)
	assert.Nil(t, err, err)

	actual, err := ungzip(tgz)
	expected := 2

	assert.Nil(t, err, err)
	assert.Equal(t, expected, actual)

	os.Remove(tgz)
}

func TestWriteFileInvalidDest(t *testing.T) {
	dest := filepath.Join("unknown", "path")
	_, err := writeFile(dest, nil)
	assert.NotNil(t, err)
}

func TestWriteFile(t *testing.T) {
	in, _ := os.CreateTemp("", "test-write-file-*")
	dest := filepath.Join(os.TempDir(), "test-dest-write-file")
	path, err := writeFile(dest, in)

	assert.Nil(t, err, err)
	assert.Equal(t, path, dest)

	os.Remove(in.Name())
	os.Remove(dest)
}

func createZipFile(dest string) error {
	file, err := os.Create(fmt.Sprintf("%s.zip", dest))
	if err != nil {
		return err
	}
	defer file.Close()

	wr := zip.NewWriter(file)
	defer wr.Close()

	_, err = wr.Create(filepath.Base(dest))
	if err != nil {
		return err
	}

	_, err = wr.Create("README.md")
	if err != nil {
		return err
	}

	return nil
}

func createGZipFile(dest string) error {
	files := []string{filepath.Base(dest), "README.md"}

	out, err := os.Create(fmt.Sprintf("%s.tar.gz", dest))
	if err != nil {
		return err
	}
	defer out.Close()

	return gzipFile(files, out)
}

func gzipFile(files []string, buf io.Writer) error {
	gw := gzip.NewWriter(buf)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, file := range files {
		err := addToArchive(tw, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func addToArchive(tw *tar.Writer, filename string) error {
	file, err := os.Create(filepath.Join(os.TempDir(), filename))
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	header.Name = filename

	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}

	os.Remove(file.Name())
	return nil
}
