package file

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Checksum(binPath, checksumsPath string) error {
	hasher := sha256.New()
	checksum, err := getChecksum(binPath, checksumsPath)
	if err != nil {
		return err
	}

	binFile, err := os.Open(binPath)
	if err != nil {
		return err
	}

	if _, err := io.Copy(hasher, binFile); err != nil {
		return err
	}

	otherChecksum := hex.EncodeToString(hasher.Sum(nil))

	if checksum != otherChecksum {
		return fmt.Errorf("checksum failed")
	}

	return nil
}

func getChecksum(binPath, checksumsPath string) (string, error) {
	file, err := os.Open(checksumsPath)
	if err != nil {
		return "", err
	}

	bytes, err := (io.ReadAll(file))
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(bytes), "\n")
	var checksum string

	for _, line := range lines {
		if strings.Contains(line, filepath.Base(binPath)) {
			columns := strings.Split(line, " ")
			checksum = columns[0]
			break
		}
	}

	return checksum, nil
}
