/*
Caravela is a Go library for updating programs.
It is intended to, given a version number, query a catalogue of releases and notify about new version or even update the program.

Usage:

	import "github.com/aureliano/caravela"

# Check for updates

Check for updates at Gitlab:

	release, err := caravela.CheckForUpdates(caravela.Conf{
		Version: "0.1.0",
		Provider: caravela.GitlabProvider{
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

Update a program taking release assets from Gitlab:

	err := caravela.Update(caravela.Conf{
		ProcessName: "oalienista",
		Version:     "0.1.0",
		Provider: caravela.GitlabProvider{
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

# Put it all together:

	release, err := caravela.CheckForUpdates(caravela.Conf{
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
			err = caravela.Update(caravela.Conf{
				ProcessName: "oalienista",
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
