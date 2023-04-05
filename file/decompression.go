package file

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Decompress(src string) (int, error) {
	if strings.HasSuffix(src, ".zip") {
		return unzip(src)
	} else if strings.HasSuffix(src, ".tar.gz") || strings.HasSuffix(src, ".tgz") {
		return ungzip(src)
	} else {
		ext := filepath.Ext(src)
		return 0, fmt.Errorf("%s não suportado para descompressão", ext)
	}
}

func unzip(src string) (int, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return 0, err
	}
	defer r.Close()

	dir := filepath.Dir(src)
	counter := 0
	for _, file := range r.File {
		in, err := file.Open()
		if err != nil {
			return counter, err
		}
		defer in.Close()

		path := filepath.Join(dir, file.Name)
		_, err = writeFile(path, in)
		if err != nil {
			return counter, err
		}

		counter++
	}

	return counter, nil
}

func ungzip(src string) (int, error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	reader, err := gzip.NewReader(srcFile)
	if err != nil {
		return 0, err
	}

	dir := filepath.Dir(src)
	return untar(dir, reader)
}

func untar(dir string, in io.Reader) (int, error) {
	tarReader := tar.NewReader(in)
	counter := 0

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return counter, err
		}

		name := header.Name

		switch header.Typeflag {
		case tar.TypeDir:
			{
				if err := os.Mkdir(name, 0755); err != nil {
					return counter, err
				}
			}
		case tar.TypeReg:
			{
				path := filepath.Join(dir, name)
				_, err = writeFile(path, tarReader)
				if err != nil {
					return counter, err
				}
			}
		default:
			return counter, fmt.Errorf("tipo do arquivo %s desconhecido (%v?)", name, header.Typeflag)
		}

		counter++
	}

	return counter, nil
}

func writeFile(dest string, in io.Reader) (string, error) {
	out, err := os.Create(dest)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(out, in)
	out.Close()

	return dest, err
}
