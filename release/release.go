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
