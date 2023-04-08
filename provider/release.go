package provider

import "time"

// A Release is a basic structured data type that abstracts a project release.
type Release struct {
	Name        string
	Description string
	ReleasedAt  time.Time
	Assets      []struct {
		Name string
		URL  string
	}
}

// Comparator is an interface which tells a type what to do to compare to releases.
type Comparator interface {
	CompareTo(r2 *Release) int
}

// CompareTo compares release r1 with release r2.
// It returns 1 if it is greater, 0 if they're equal or -1 likewise.
func (r1 *Release) CompareTo(r2 *Release) int {
	return compareVersions(r1.Name, r2.Name)
}
