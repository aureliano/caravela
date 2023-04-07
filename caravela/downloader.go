package caravela

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	httpc "github.com/aureliano/caravela/http"
)

func downloadTo(client httpc.HttpClientPlugin, release *Release, dir string) (string, string, error) {
	fname, furl := findReleaseFileUrl(runtime.GOOS, release)
	if fname == "" {
		return "", "", fmt.Errorf("não há uma versão compatível com %s", runtime.GOOS)
	}

	fileBin := filepath.Join(dir, fname)

	wmsg(100)
	err := downloadFile(client, furl, fileBin)
	if err != nil {
		return "", "", err
	}

	fname = "checksums.txt"
	furl = findChecksumsFileUrl(release)
	fileChecksums := filepath.Join(dir, fname)
	if furl == "" {
		return "", "", fmt.Errorf("não encontrou o arquivo %s", fname)
	}

	wmsg(101)
	err = downloadFile(client, furl, fileChecksums)
	if err != nil {
		return "", "", err
	}

	return fileBin, fileChecksums, nil
}

func downloadFile(client httpc.HttpClientPlugin, sourceUrl, dest string) error {
	file, err := os.Create(dest)
	if err != nil {
		os.Remove(dest)
		return err
	}
	defer file.Close()

	req, _ := http.NewRequest(http.MethodGet, sourceUrl, nil)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro na integração com o Gitlab: %d", resp.StatusCode)
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func findReleaseFileUrl(osys string, release *Release) (string, string) {
	for _, asset := range release.Assets {
		name := strings.ToLower(asset.Name)
		if strings.Contains(name, osys) {
			return asset.Name, asset.URL
		}
	}

	return "", ""
}

func findChecksumsFileUrl(release *Release) string {
	for _, asset := range release.Assets {
		if asset.Name == "checksums.txt" {
			return asset.URL
		}
	}

	return ""
}
