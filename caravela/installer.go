package caravela

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var osExecutable = os.Executable

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
		if strings.HasSuffix(file.Name(), "tar.gz") || strings.HasSuffix(file.Name(), "zip") || file.Name() == "checksums.txt" {
			continue
		}

		src := filepath.Join(srcDir, file.Name())
		dest := filepath.Join(dir, file.Name())

		wmsg(300, src, dest)
		err = installFile(dest, src)
		if err != nil {
			return nil
		}
	}

	return nil
}

func homeDir() (string, error) {
	ex, err := osExecutable()

	if err != nil {
		return "", fmt.Errorf("erro ao obter binário em execução: %w", err)
	}

	fi, err := os.Lstat(ex)
	if err != nil {
		return "", fmt.Errorf("erro ao obter informacões sobre o binário: %w", err)
	}

	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		path, err := filepath.EvalSymlinks(ex)
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
		fm = fs.FileMode(0644)
	} else if err != nil {
		return err
	}

	if fileExist {
		if err = os.Remove(dest); err != nil {
			return err
		}
		fm = fs.FileMode(destInfo.Mode())
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
