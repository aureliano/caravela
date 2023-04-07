package caravela

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func serializeRelease(release *Release) error {
	source, err := json.Marshal(release)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	fname := fmt.Sprintf("release_%s.json", now.Format("2006-01-02"))

	file := filepath.Join(os.TempDir(), fname)
	err = os.WriteFile(file, source, 0644)

	return err
}

func deserializeRelease() (*Release, error) {
	now := time.Now().UTC()
	fname := fmt.Sprintf("release_%s.json", now.Format("2006-01-02"))

	file := filepath.Join(os.TempDir(), fname)
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	release := &Release{}
	err = json.Unmarshal(bytes, release)

	return release, err
}
