package caravela

type UpdaterProvider interface {
	FetchLastRelease(client HttpClientPlugin) (*Release, error)
	CacheRelease(release Release) error
	RestoreCacheRelease() (*Release, error)
}
