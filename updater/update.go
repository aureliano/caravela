package updater

import (
	"os"
	"path/filepath"

	pvdr "github.com/aureliano/caravela/provider"
)

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
	pname,
	currver string,
) (*pvdr.Release, error) {
	rel, err := FindUpdate(client, provider, currver)
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(os.TempDir(), pname)
	_ = os.MkdirAll(dir, os.ModePerm)

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

	err = mpInstall(dir)
	if err != nil {
		return nil, err
	}

	os.RemoveAll(dir)

	return rel, nil
}
