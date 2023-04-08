package main

import (
	"fmt"

	"github.com/aureliano/caravela/caravela"
	"github.com/aureliano/caravela/provider"
)

func main() {
	err := caravela.Update(caravela.Conf{
		ProcessName: "bruzundangas",
		Version:     "0.1.0",
		Provider: provider.GitlabProvider{
			Host:        "gitlab.com",
			Ssl:         true,
			ProjectPath: "gitlab-org/gitlab",
		},
	})

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Process successfuly updated!")
	}
}
