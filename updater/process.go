package updater

import (
	"fmt"
	"os"
	"path/filepath"
)

var _mpOsExecutable = os.Executable

func processFilePath() (string, error) {
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

		return path, nil
	}

	return ex, nil
}
