package main

import (
	"log"

	"github.com/aureliano/caravela"
	"github.com/aureliano/caravela/provider"
)

func main() {
	release, err := caravela.CheckUpdates(caravela.Conf{
		Version:     "v1.0.0",
		IgnoreCache: true,
		Provider: provider.GitlabProvider{
			Host:        "gitlab.com",
			Ssl:         true,
			ProjectPath: "opennota/fb2index",
		},
	})

	if err != nil {
		log.Fatalln(err)
	}

	expected := "v1.0.3"

	if expected != release.Name {
		log.Fatalf("Expected %s, but got %s instead.", expected, release.Name)
	}
}
