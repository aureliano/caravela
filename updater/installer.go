package updater

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var _mpOsExecutable = os.Executable

func install(srcDir string) error {
	dir, err := homeDir()
	if err != nil {
		return err
	}

	files, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), "tar.gz") ||
			strings.HasSuffix(file.Name(), "zip") ||
			file.Name() == "checksums.txt" {
			continue
		}

		src := filepath.Join(srcDir, file.Name())
		dest := filepath.Join(dir, file.Name())

		err = installFile(dest, src)
		if err != nil {
			return err
		}
	}

	return nil
}

func homeDir() (string, error) {
	ex, err := _mpOsExecutable()

	if err != nil {
		return "", fmt.Errorf("error getting running process: %w", err)
	}

	fi, err := os.Lstat(ex)
	if err != nil {
		return "", fmt.Errorf("error getting information from process: %w", err)
	}

	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		var path string
		path, err = filepath.EvalSymlinks(ex)

		if err != nil {
			return "", err
		}

		return filepath.Dir(path), nil
	}

	return filepath.Dir(ex), nil
}

func installFile(dest, src string) error {
	destInfo, err := os.Stat(dest)
	fileExist := true
	var fm fs.FileMode

	if os.IsNotExist(err) {
		fileExist = false
		const permFile = 0644
		fm = fs.FileMode(permFile)
	} else if err != nil {
		return err
	}

	if fileExist {
		if err = os.Remove(dest); err != nil {
			return err
		}
		fm = destInfo.Mode()
	}

	out, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fm)
	if err != nil {
		return err
	}
	defer out.Close()

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return nil
}
