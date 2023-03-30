package release

import "time"

type Release struct {
	Name        string
	Description string
	ReleasedAt  time.Time
	Assets      []struct {
		Name string
		URL  string
	}
}

type Comparator interface {
	CompareTo(r2 *Release) int
}

func (r1 *Release) CompareTo(r2 *Release) int {
	return CompareVersions(r1.Name, r2.Name)
}
