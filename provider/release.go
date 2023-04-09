package provider

import "time"

// Release is a structured data type that abstracts the release of a project.
type Release struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ReleasedAt  time.Time `json:"releasedAt"`
	Assets      []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"assets"`
}

// Comparator is an interface that provides methods for comparing two releases.
type Comparator interface {
	CompareTo(r2 *Release) int
}

// CompareTo compares release r1 with release r2.
// It returns 1, 0 or -1 if it is greater, equal or lesser.
func (r1 *Release) CompareTo(r2 *Release) int {
	return compareVersions(r1.Name, r2.Name)
}
