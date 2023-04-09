package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// GitlabProvider is a provider for getting releases from Gitlab.
type GitlabProvider struct {
	Host        string
	Port        uint
	Ssl         bool
	ProjectPath string
	Timeout     time.Duration
}

// GitlabRelease is a representation - in JSON form - of what Gitlab
// returns when the target service is called.
type GitlabRelease struct {
	Name        string    `json:"tag_name"`
	Description string    `json:"description"`
	ReleaseAt   time.Time `json:"released_at"`
	Assets      struct {
		Links []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"links"`
	} `json:"assets"`
}

func (provider GitlabProvider) FetchLastRelease(client HTTPClientPlugin) (*Release, error) {
	initProvider(&provider)
	err := validateProvider(provider)
	if err != nil {
		return nil, err
	}

	releases, err := fetchReleases(provider, client)
	if err != nil {
		return nil, err
	}

	var lastRelease *Release
	for _, release := range releases {
		if lastRelease == nil {
			lastRelease = release
		} else if lastRelease.CompareTo(release) == -1 {
			lastRelease = release
		}
	}

	return lastRelease, nil
}

func (GitlabProvider) CacheRelease(r Release) error {
	return serializeRelease(&r)
}

func (GitlabProvider) RestoreCacheRelease() (*Release, error) {
	return deserializeRelease()
}

func (r1 *GitlabRelease) CompareTo(r2 *GitlabRelease) int {
	return compareVersions(r1.Name, r2.Name)
}

func fetchReleases(p GitlabProvider, client HTTPClientPlugin) ([]*Release, error) {
	srvURL := buildServiceURL(p)
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, srvURL, nil)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gitlab integration error: %d", resp.StatusCode)
	}

	var releases []*GitlabRelease
	err = json.NewDecoder(resp.Body).Decode(&releases)

	return convertReleases(releases), err
}

func buildServiceURL(p GitlabProvider) string {
	projectPath := url.QueryEscape(p.ProjectPath)
	protocol := "http"
	if p.Ssl {
		protocol += "s"
	}
	baseURL := fmt.Sprintf("%s://%s:%d/api/v4/projects", protocol, p.Host, p.Port)

	return fmt.Sprintf("%s/%s/releases", baseURL, projectPath)
}

func convertReleases(in []*GitlabRelease) []*Release {
	size := len(in)
	rels := make([]*Release, size)

	for i, r := range in {
		rels[i] = convertToBase(r)
	}

	return rels
}

func convertToBase(r *GitlabRelease) *Release {
	t := Release{
		Name:        r.Name,
		Description: r.Description,
		ReleasedAt:  r.ReleaseAt,
	}

	size := len(r.Assets.Links)
	t.Assets = make([]struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}, size)

	for i, link := range r.Assets.Links {
		t.Assets[i] = struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}{Name: link.Name, URL: link.URL}
	}

	return &t
}

func validateProvider(p GitlabProvider) error {
	switch {
	case p.Host == "":
		return fmt.Errorf("host is required")
	case p.Port <= 0:
		return fmt.Errorf("port must be > 0")
	case p.ProjectPath == "":
		return fmt.Errorf("project path is required")
	default:
		return nil
	}
}

func initProvider(p *GitlabProvider) {
	const httpPort = 80
	const httpsPort = 443
	const timeout = time.Second * 30

	if p.Port == 0 {
		if p.Ssl {
			p.Port = httpsPort
		} else {
			p.Port = httpPort
		}
	} else {
		p.Ssl = p.Port == httpsPort
	}

	if p.Timeout == 0 {
		p.Timeout = timeout
	}
}
