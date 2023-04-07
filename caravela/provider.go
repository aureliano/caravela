package caravela

type UpdaterProvider interface {
	FetchLastRelease(client httpClientPlugin) (*Release, error)
	CacheRelease(release Release) error
	RestoreCacheRelease() (*Release, error)
}
