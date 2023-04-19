/*
Caravela is a Go library to support program update automation.

Some platforms, such as GitHub and GitLab, provide an API for querying and retrieving software versions.
Indeed, this library queries this API to check for new versions and even updates the program.

GitHub releases: https://docs.github.com/en/rest/releases
GitLab releases: https://docs.gitlab.com/ee/api/releases

Usage:

	import "github.com/aureliano/caravela"

# Check updates

CheckUpdates fetches the last release published.
It returns the last release available or raises an error if the current version is already the last one.

	release, err := caravela.CheckUpdates(caravela.Conf{
		Version: "0.1.0",
		Provider: provider.GitlabProvider{
			Host:        "gitlab.com",
			Ssl:         true,
			ProjectPath: "gitlab-org/gitlab",
		},
	})

	if err != nil {
		fmt.Printf("Check for updates has failed! %s\n", err)
	} else {
		fmt.Printf("New version available %s\n%s\n", release.Name, release.Description)
	}

# Update

Update updates running program to the last available release.
It returns the release used to update this program or raises an error if it's already the last version.

	release, err := caravela.Update(caravela.Conf{
		Version:     "0.1.0",
		Provider: provider.GitlabProvider{
			Host:        "gitlab.com",
			Ssl:         true,
			ProjectPath: "gitlab-org/gitlab",
		},
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("New version installed!")
	}

# Put it all together

Let's put it all together chainning CheckUpdates and Update.

	release, err := caravela.CheckUpdates(caravela.Conf{
		Version: "0.1.0",
		Provider: provider.GitlabProvider{
			Host:        "gitlab.com",
			Ssl:         true,
			ProjectPath: "gitlab-org/gitlab",
		},
	})

	if err != nil {
		fmt.Printf("Check for updates has failed! %s\n", err)
	} else {
		fmt.Printf("[WARN] There is a new version available. Would you like to update this program?")

		// ...
		// Ask user whether to update or not.
		// ...

		if shouldUpdate {
			release, err := caravela.Update(caravela.Conf{
				Version:     "0.1.0",
				Provider: provider.GitlabProvider{
					Host:        "gitlab.com",
					Ssl:         true,
					ProjectPath: "gitlab-org/gitlab",
				},
			})

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			} else {
				fmt.Printf("New version %s was successfuly installed!\n", release.Name)
			}
		}
	}
*/
package caravela
