package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/aureliano/caravela"
	"github.com/aureliano/caravela/provider"
)

func main() {
	release, err := caravela.CheckUpdates(caravela.Conf{
		Version:     "v0.1.0-alpha",
		IgnoreCache: true,
		Provider: provider.GithubProvider{
			Host:        "api.github.com",
			Ssl:         true,
			ProjectPath: "aureliano/caravela",
		},
	})

	if err != nil {
		log.Fatalln(err)
	}

	tag, err := lastTaggedVersion()
	if err != nil {
		log.Fatalln(err)
	}

	if tag != release.Name {
		log.Fatalf("Expected %s, but got %s instead.", tag, release.Name)
	}
}

func lastTaggedVersion() (string, error) {
	cmd := exec.Command("git", "tag", "-l")
	out, err := cmd.Output()
	strout := string(out)

	if err != nil {
		return "", err
	}

	tags := strings.Split(strings.Trim(strout, "\n"), "\n")
	if len(tags) == 0 {
		return "", fmt.Errorf("no tag to compare with")
	}

	return tags[len(tags)-1], nil
}
