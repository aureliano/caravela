package updater

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func decompress(src string) (int, error) {
	switch {
	case strings.HasSuffix(src, ".zip"):
		return unzip(src)
	case strings.HasSuffix(src, ".tar.gz") || strings.HasSuffix(src, ".tgz"):
		return ungzip(src)
	default:
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
		_, e = writeFile(path, in, file.Mode())
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
				cname := filepath.Clean(name)
				path := filepath.Join(dir, cname)
				virtualDir := filepath.Dir(cname)
				if virtualDir != "." {
					err = os.MkdirAll(filepath.Join(dir, virtualDir), fs.ModePerm)
					if err != nil {
						return counter, err
					}
				}

				_, err = writeFile(path, tarReader, fs.FileMode(header.Mode))
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

func writeFile(dest string, in io.Reader, perm fs.FileMode) (string, error) {
	out, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(out, in)
	out.Close()

	return dest, err
}
