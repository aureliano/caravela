package updater

import (
	"fmt"
	"os"
	"path/filepath"

	pvdr "github.com/aureliano/caravela/provider"
)

var mpProcessFilePath = processFilePath
var mpDownloadTo = downloadTo
var mpDecompress = decompress
var mpChecksum = checksum
var mpInstall = install

// Update updates running program to the last available release.
//
// It returns the release used to update this program or raises
// an error if it's already the last version.
func UpdateRelease(
	client pvdr.HTTPClientPlugin,
	provider pvdr.UpdaterProvider,
	currver string,
	ignoreCache bool,
) (*pvdr.Release, error) {
	rel, err := FindUpdate(client, provider, currver, ignoreCache)
	if err != nil {
		return nil, err
	} else if rel.Name == "" {
		return nil, fmt.Errorf("already on the edge")
	}

	procFile, err := mpProcessFilePath()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(os.TempDir(), filepath.Base(procFile))
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	bin, checksums, err := mpDownloadTo(client, rel, dir)
	if err != nil {
		return nil, err
	}

	_, err = mpDecompress(bin)
	if err != nil {
		return nil, err
	}

	err = mpChecksum(bin, checksums)
	if err != nil {
		return nil, err
	}

	err = mpInstall(dir, filepath.Dir(procFile))
	if err != nil {
		return nil, err
	}

	os.RemoveAll(dir)

	return rel, nil
}
