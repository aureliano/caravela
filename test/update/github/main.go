package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aureliano/caravela"
	"github.com/aureliano/caravela/provider"
)

func main() {
	release, err := caravela.Update(caravela.Conf{
		Version:     "v1.0.0",
		IgnoreCache: true,
		Provider: provider.GithubProvider{
			Host:        "api.github.com",
			Ssl:         true,
			ProjectPath: "goreleaser/goreleaser",
		},
	})

	if err != nil {
		log.Fatalln(err)
	}

	expected := "v1.17.1"

	if release.Name != expected {
		log.Fatalf("Expected %s, but got %s instead.", expected, release.Name)
	}

	err = assertAssetsArePresent()
	if err != nil {
		log.Fatalln(err)
	}
}

func assertAssetsArePresent() error {
	exec, err := os.Executable()
	if err != nil {
		return err
	}

	dir := filepath.Dir(exec)
	err = filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if path == exec {
				return nil
			}

			if !info.IsDir() && !assetIsPresent(path) {
				return fmt.Errorf("unexpected file %s", path)
			}

			return nil
		})

	return err
}

func assetIsPresent(path string) bool {
	assets := []string{
		"LICENSE.md",
		"README.md",
		"completions/goreleaser.bash",
		"completions/goreleaser.fish",
		"completions/goreleaser.zsh",
		"goreleaser",
		"manpages/goreleaser.1.gz",
	}

	for _, asset := range assets {
		if strings.HasSuffix(path, asset) {
			return true
		}
	}

	return false
}
