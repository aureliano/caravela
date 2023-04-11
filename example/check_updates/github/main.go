package main

import (
	"fmt"

	"github.com/aureliano/caravela"
	"github.com/aureliano/caravela/provider"
)

func main() {
	release, err := caravela.CheckUpdates(caravela.Conf{
		Version: "0.1.0",
		Provider: provider.GithubProvider{
			Host:        "api.github.com",
			Ssl:         true,
			ProjectPath: "aureliano/caravela",
		},
	})

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Version: %s\n", release.Name)
		fmt.Printf("Description: %s\n", release.Description)
		fmt.Printf("Date release: %v\n", release.ReleasedAt)
		fmt.Printf("Assets: %v\n", release.Assets)
		fmt.Println("github")
	}
}
