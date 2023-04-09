/*
The provider package contains the release recovery logic of the various project hosting systems.
For each system, there is a file that implements the UpdaterProvider interface.
Thus, for GitHub and GitLab systems, we will have github.go and gitlab.go files, respectively.

# UpdaterProvider

	// It is the interface that every release provider should implement, as it has all the
	// expected method definitions for querying, caching and restoring cached release.
	type UpdaterProvider interface {
		// FetchLastRelease queries provider for the last release of a project.
		FetchLastRelease(client HTTPClientPlugin) (*Release, error)

		// CacheRelease writes the release passed as parameter to the file system.
		CacheRelease(release Release) error

		// RestoreCacheRelease retores a cached release.
		// It returns nil if no release has been cached yet.
		RestoreCacheRelease() (*Release, error)
	}

# GitLab provider implementation

	// GitlabProvider is a provider for getting releases from Gitlab.
	type GitlabProvider struct {
		Host        string
		Port        uint
		Ssl         bool
		ProjectPath string
		Timeout     time.Duration
	}

	func (provider GitlabProvider) FetchLastRelease(client HTTPClientPlugin) (*Release, error) {
		// ...
	}

	func (GitlabProvider) CacheRelease(r Release) error {
		return serializeRelease(&r)
	}

	func (GitlabProvider) RestoreCacheRelease() (*Release, error) {
		return deserializeRelease()
	}
*/
package provider
