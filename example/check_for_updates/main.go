package main

import (
	"fmt"

	"github.com/aureliano/caravela/caravela"
)

func main() {
	release, err := caravela.CheckForUpdates(caravela.Conf{
		Version: "0.1.0",
		Provider: caravela.GitlabProvider{
			Host:        "gitlab.com",
			Ssl:         true,
			ProjectPath: "gitlab-org/gitlab",
		},
	})

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Version: %s\n", release.Name)
		fmt.Printf("Description: %s\n", release.Description)
		fmt.Printf("Date release: %v\n", release.ReleasedAt)
		fmt.Printf("Assets: %v\n", release.Assets)
	}
}
