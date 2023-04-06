/*
Caravela is a Go library for updating programs.
It is intended to, given a version number, query a catalogue of releases and notify about new version or even update the program.

Usage:

	import "github.com/aureliano/caravela"

# Check for updates
Check for updates at Gitlab:

	client, _ := http.BuildClientTls12()
	releaseProvider := provider.GitlabProvider{
		Host:        "gitlab.com",
		Port:        80,
		Ssl:         true,
		ProjectPath: "massis/oalienista",
	}
	conf := i18n.I18nConf{Verbose: true, Locale: i18n.EN}
	release, err := caravela.CheckForUpdates(client, releaseProvider, conf, "0.1.0")

	if err != nil {
		fmt.Printf("Check for updates has failed! %s\n", err)
	} else {
		fmt.Printf("New version available %s\n%s\n", release.Name, release.Description)
	}

# Update
Update a program taking release assets from Gitlab:

	client, _ := http.BuildClientTls12()
	releaseProvider := provider.GitlabProvider{
		Host:        "gitlab.com",
		Port:        80,
		Ssl:         true,
		ProjectPath: "massis/oalienista",
	}
	conf := i18n.I18nConf{Verbose: true, Locale: i18n.EN}
	err := caravela.Update(client, releaseProvider, conf, "oalienista", "0.1.0")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("New version installed!")
	}

# Put it all together:

	client, _ := http.BuildClientTls12()
	releaseProvider := provider.GitlabProvider{
		Host:        "gitlab.com",
		Port:        80,
		Ssl:         true,
		ProjectPath: "massis/oalienista",
	}
	conf := i18n.I18nConf{Verbose: true, Locale: i18n.EN}
	release, err := caravela.CheckForUpdates(client, releaseProvider, conf, "0.1.0")

	if err != nil {
		fmt.Printf("Check for updates has failed! %s\n", err)
	} else {
		fmt.Printf("[WARN] There is a new version available. Would you like to update this program?")

		// ...
		// Ask user whether to update or not.
		// ...

		if shouldUpdate {
			err = caravela.Update(c, p, conf, "oalienista", "0.1.0")
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
