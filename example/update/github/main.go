package main

import (
	"fmt"

	"github.com/aureliano/caravela"
	"github.com/aureliano/caravela/provider"
)

func main() {
	release, err := caravela.Update(caravela.Conf{
		ProcessName: "bruzundangas",
		Version:     "0.1.0",
		IgnoreCache: true,
		Provider: provider.GithubProvider{
			Host:        "api.github.com",
			Ssl:         true,
			ProjectPath: "aureliano/caravela",
		},
	})

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Process successfuly updated!")
		fmt.Printf("%s: %s", release.Name, release.Description)
	}
}
