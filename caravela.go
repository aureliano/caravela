package caravela

import (
	"fmt"
	"net/http"

	pvdr "github.com/aureliano/caravela/provider"
	caravela "github.com/aureliano/caravela/updater"
)

// A Conf is a wrapper of data to be passed as input to the public functions.
type Conf struct {
	ProcessName string
	Version     string
	Provider    pvdr.UpdaterProvider
	HTTPClient  *http.Client
}

var mpCheckForUpdates = caravela.FindUpdate
var mpUpdate = caravela.UpdateRelease

// CheckUpdates fetches the last release published.
//
// It returns the last release available or raises an error
// if the current version is already the last one.
func CheckUpdates(c Conf) (*pvdr.Release, error) {
	if c.Version == "" {
		return nil, fmt.Errorf("current version is required")
	}

	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}

	client := pvdr.HTTPClientDecorator{Client: *c.HTTPClient}

	return mpCheckForUpdates(&client, c.Provider, c.Version)
}

// Update updates running program to the last available release.
//
// It returns the release used to update this program or raises
// an error if it's already the last version.
func Update(c Conf) (*pvdr.Release, error) {
	if c.ProcessName == "" {
		return nil, fmt.Errorf("process name is required")
	}

	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}

	client := pvdr.HTTPClientDecorator{Client: *c.HTTPClient}

	return mpUpdate(&client, c.Provider, c.ProcessName, c.Version)
}
