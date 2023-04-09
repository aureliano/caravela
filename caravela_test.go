package caravela

import (
	"fmt"
	"net/http"
	"testing"

	pvdr "github.com/aureliano/caravela/provider"
	"github.com/stretchr/testify/assert"
)

func TestCheckForUpdatesCurrentVersionIsRequired(t *testing.T) {
	_, err := CheckUpdates(Conf{})
	assert.Equal(t, "current version is required", err.Error())
}

func TestCheckForUpdatesHTTPClientIsNil(t *testing.T) {
	mpCheckForUpdates = func(client pvdr.HTTPClientPlugin, provider pvdr.UpdaterProvider,
		currver string) (*pvdr.Release, error) {
		return nil, fmt.Errorf("already on the edge")
	}

	r, err := CheckUpdates(Conf{Version: "0.1.0"})
	assert.Nil(t, r)
	assert.Equal(t, "already on the edge", err.Error())
}

func TestCheckForUpdates(t *testing.T) {
	mpCheckForUpdates = func(client pvdr.HTTPClientPlugin, provider pvdr.UpdaterProvider,
		currver string) (*pvdr.Release, error) {
		return nil, fmt.Errorf("already on the edge")
	}

	r, err := CheckUpdates(Conf{HTTPClient: http.DefaultClient, Version: "0.1.0"})
	assert.Nil(t, r)
	assert.Equal(t, "already on the edge", err.Error())
}

func TestUpdateProcessNameIsRequired(t *testing.T) {
	_, err := Update(Conf{ProcessName: "", Version: "0.1.0"})
	assert.Equal(t, "process name is required", err.Error())
}

func TestUpdateHTTPClientIsNil(t *testing.T) {
	mpUpdate = func(client pvdr.HTTPClientPlugin, provider pvdr.UpdaterProvider,
		pname, currver string) (*pvdr.Release, error) {
		return nil, fmt.Errorf("")
	}

	_, err := Update(Conf{ProcessName: "oalienista", Version: "0.1.0"})
	assert.NotNil(t, err)
}
