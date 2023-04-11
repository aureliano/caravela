package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GithubProvider is a provider for getting releases from Github.
type GithubProvider struct {
	Host        string
	Port        uint
	Ssl         bool
	ProjectPath string
	Timeout     time.Duration
}

// GithubRelease is a representation - in JSON form - of what Github
// returns when the target service is called.
type GithubRelease struct {
	Name        string    `json:"tag_name"`
	Body        string    `json:"body"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	} `json:"assets"`
}

func (provider GithubProvider) FetchLastRelease(client HTTPClientPlugin) (*Release, error) {
	initGithubProvider(&provider)
	err := validateGithubProvider(provider)
	if err != nil {
		return nil, err
	}

	releases, err := fetchGithubReleases(provider, client)
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

func (GithubProvider) CacheRelease(r Release) error {
	return serializeRelease(&r)
}

func (GithubProvider) RestoreCacheRelease() (*Release, error) {
	return deserializeRelease()
}

func (r1 *GithubRelease) CompareTo(r2 *GithubRelease) int {
	return compareVersions(r1.Name, r2.Name)
}

func fetchGithubReleases(p GithubProvider, client HTTPClientPlugin) ([]*Release, error) {
	srvURL := buildGithubServiceURL(p)
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, srvURL, nil)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github integration error: %d", resp.StatusCode)
	}

	var releases []*GithubRelease
	err = json.NewDecoder(resp.Body).Decode(&releases)

	return convertGithubReleases(releases), err
}

func buildGithubServiceURL(p GithubProvider) string {
	protocol := "http"
	if p.Ssl {
		protocol += "s"
	}
	baseURL := fmt.Sprintf("%s://%s:%d/repos", protocol, p.Host, p.Port)

	return fmt.Sprintf("%s/%s/releases", baseURL, p.ProjectPath)
}

func convertGithubReleases(in []*GithubRelease) []*Release {
	size := len(in)
	rels := make([]*Release, size)

	for i, r := range in {
		rels[i] = convertGithubToBase(r)
	}

	return rels
}

func convertGithubToBase(r *GithubRelease) *Release {
	t := Release{
		Name:        r.Name,
		Description: r.Body,
		ReleasedAt:  r.PublishedAt,
	}

	size := len(r.Assets)
	t.Assets = make([]struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}, size)

	for i, link := range r.Assets {
		t.Assets[i] = struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}{Name: link.Name, URL: link.URL}
	}

	return &t
}

func validateGithubProvider(p GithubProvider) error {
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

func initGithubProvider(p *GithubProvider) {
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
