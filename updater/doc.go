/*
The updater package has the logic of the functions that are exported by this library.
In short, calls made to functions exported to the client are delegated to this package.
Indeed, you can even make direct calls to this package, but having to fill in all the input parameters.

Usage:

	import "github.com/aureliano/caravela/updater"

# Find update

Fetches the last release published.
It returns the last release available or raises an error if the current version is already the last one.

	release, err := updater.FindUpdate(
		&provider.HTTPClientDecorator{Client: *http.DefaultClient},
		provider.GitlabProvider{
			Host:        "gitlab.com",
			Ssl:         true,
			ProjectPath: "gitlab-org/gitlab",
		},
		"0.1.0",
	)

# Update

Updates running program to the last available release.
It returns the release used to update this program or raises an error if it's already the last version.

	release, err := updater.UpdateRelease(
		&provider.HTTPClientDecorator{Client: *http.DefaultClient},
		provider.GitlabProvider{
			Host:        "gitlab.com",
			Ssl:         true,
			ProjectPath: "gitlab-org/gitlab",
		},
		"oalienista",
		"0.1.0",
	)
*/
package updater
