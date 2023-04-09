package caravela

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/aureliano/caravela/provider"
)

const downloadTimeout = time.Second * 120
const checksumsFileName = "checksums.txt"

var mpDownloadFile = downloadFile

func downloadTo(client provider.HTTPClientPlugin, release *provider.Release, dir string) (string, string, error) {
	fname, furl := findReleaseFileURL(runtime.GOOS, release)
	if fname == "" {
		return "", "", fmt.Errorf("there is no version compatible with %s", runtime.GOOS)
	}

	fileBin := filepath.Join(dir, fname)

	err := mpDownloadFile(client, furl, fileBin)
	if err != nil {
		return "", "", err
	}

	furl = findChecksumsFileURL(release)
	fileChecksums := filepath.Join(dir, checksumsFileName)
	if furl == "" {
		return "", "", fmt.Errorf("file %s not found", checksumsFileName)
	}

	err = mpDownloadFile(client, furl, fileChecksums)
	if err != nil {
		return "", "", err
	}

	return fileBin, fileChecksums, nil
}

func downloadFile(client provider.HTTPClientPlugin, sourceURL, dest string) error {
	file, err := os.Create(dest)
	if err != nil {
		os.Remove(dest)
		return err
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), downloadTimeout)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http error (%d)", resp.StatusCode)
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func findReleaseFileURL(osys string, release *provider.Release) (string, string) {
	for _, asset := range release.Assets {
		name := strings.ToLower(asset.Name)
		if strings.Contains(name, osys) {
			return asset.Name, asset.URL
		}
	}

	return "", ""
}

func findChecksumsFileURL(release *provider.Release) string {
	for _, asset := range release.Assets {
		if asset.Name == checksumsFileName {
			return asset.URL
		}
	}

	return ""
}
