package main

import (
	"log"

	"github.com/aureliano/caravela"
	"github.com/aureliano/caravela/provider"
)

func main() {
	_, err := caravela.Update(caravela.Conf{
		Version:     "v1.0.0",
		IgnoreCache: true,
		Provider: provider.GitlabProvider{
			Host:        "gitlab.com",
			Ssl:         true,
			ProjectPath: "commonground/haven/haven",
		},
	})

	expected := "file checksums.txt not found"
	if err == nil {
		log.Fatalf("Expected error: %s\n", expected)
	} else if expected != err.Error() {
		log.Fatalf("Expected error '%s', but got '%s' instead.", expected, err.Error())
	}
}
