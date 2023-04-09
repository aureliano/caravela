package caravela

import (
	"fmt"
	"net/http"

	pvdr "github.com/aureliano/caravela/provider"
)

// A Conf is a wrapper o data to be passed as input to the public functions.
type Conf struct {
	ProcessName string
	Version     string
	Provider    pvdr.UpdaterProvider
	HTTPClient  *http.Client
}

var mpCheckForUpdates = FindUpdate
var mpUpdate = UpdateRelease

// CheckForUpdates queries, given a provider, for new releases.
// It returns the last release available or nil if the current
// version is already the last one.
func CheckForUpdates(c Conf) (*pvdr.Release, error) {
	if c.Version == "" {
		return nil, fmt.Errorf("current version is required")
	}

	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}

	client := pvdr.HTTPClientDecorator{Client: *c.HTTPClient}

	return mpCheckForUpdates(&client, c.Provider, c.Version)
}

// Update running program to the last available release.
// Raises an error if it's already the last version
// or returns the new release.
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
