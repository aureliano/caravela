package caravela

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func decompress(src string) (int, error) {
	if strings.HasSuffix(src, ".zip") {
		return unzip(src)
	} else if strings.HasSuffix(src, ".tar.gz") || strings.HasSuffix(src, ".tgz") {
		return ungzip(src)
	} else {
		ext := filepath.Ext(src)
		return 0, fmt.Errorf("%s not supported for decompression", ext)
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
		in, e := file.Open()
		if e != nil {
			return counter, e
		}
		defer in.Close()

		if strings.Contains(dir, "..") || strings.Contains(file.Name, "..") {
			return counter, fmt.Errorf("err")
		}

		path := filepath.Join(dir, filepath.Clean(file.Name))
		_, e = writeFile(path, in)
		if e != nil {
			return counter, e
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

		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return counter, err
		}

		name := header.Name

		switch header.Typeflag {
		case tar.TypeDir:
			{
				if err = os.Mkdir(name, 0755); err != nil {
					return counter, err
				}
			}
		case tar.TypeReg:
			{
				path := filepath.Join(dir, filepath.Clean(name))
				_, err = writeFile(path, tarReader)
				if err != nil {
					return counter, err
				}
			}
		default:
			return counter, fmt.Errorf("%s has unknown type (%v?)", name, header.Typeflag)
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
