package caravela

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type GitlabProvider struct {
	Host        string
	Port        uint
	Ssl         bool
	ProjectPath string
}

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

func (provider GitlabProvider) FetchLastRelease(client httpClientPlugin) (*Release, error) {
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
		} else {
			if lastRelease.CompareTo(release) == -1 {
				lastRelease = release
			}
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

func fetchReleases(p GitlabProvider, client httpClientPlugin) ([]*Release, error) {
	srvUrl := buildServiceUrl(p)
	req, _ := http.NewRequest(http.MethodGet, srvUrl, nil)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Gitlab integration error: %d", resp.StatusCode)
	}

	var releases []*GitlabRelease
	err = json.NewDecoder(resp.Body).Decode(&releases)

	return convertReleases(releases), err
}

func buildServiceUrl(p GitlabProvider) string {
	projectPath := url.QueryEscape(p.ProjectPath)
	protocol := "http"
	if p.Ssl {
		protocol += "s"
	}
	baseUrl := fmt.Sprintf("%s://%s:%d/api/v4/projects", protocol, p.Host, p.Port)

	return fmt.Sprintf("%s/%s/releases", baseUrl, projectPath)
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
		Name string
		URL  string
	}, size)

	for i, link := range r.Assets.Links {
		t.Assets[i] = struct {
			Name string
			URL  string
		}{Name: link.Name, URL: link.URL}
	}

	return &t
}

func validateProvider(p GitlabProvider) error {
	if p.Host == "" {
		return fmt.Errorf("host is required")
	} else if p.Port <= 0 {
		return fmt.Errorf("port must be > 0")
	} else if p.ProjectPath == "" {
		return fmt.Errorf("project path is required")
	} else {
		return nil
	}
}
