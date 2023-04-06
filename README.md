# caravela

[![CI Pipeline](https://github.com/aureliano/caravela/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/aureliano/caravela/actions/workflows/build.yml?query=branch%3Amain)
[![Coverage](https://coveralls.io/repos/github/aureliano/caravela/badge.svg?branch=main)](https://coveralls.io/github/aureliano/caravela?branch=main)

Caravela is a Go library for updating programs.

It is intended to, given a version number, query a [catalogue of releases](https://github.com/aureliano/caravela/releases) and notify about new version or even update the program.

Currently, **caravela requires Go version 1.17 or greater**. Caravela tracks Go's version support policy. We do our best not to break older versions of Go if we don't have to, but due to tooling constraints, we don't always test older versions.

## Installation
To install Caravela, use `go get`:

`go get github.com/aureliano/caravela`

Or you can install specific version like:

`go get github.com/aureliano/caravela/v0`

Or even add it as a project depency of your module:

`require github.com/aureliano/caravela v0`

### Staying up to date
To update Caravela to the latest version, use `go get -u github.com/aureliano/caravela`.

## Usage
You'll usually call `Update` after the given result of `CheckForUpdates`. The sample bellow show you how to notice the user about new release available and ask him about updating.

```go
package yours

import (
	"fmt"
	"os"

	"github.com/aureliano/caravela"
	"github.com/aureliano/caravela/http"
	"github.com/aureliano/caravela/i18n"
	"github.com/aureliano/caravela/provider"
)

func main() {
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
	} else if release != nil {
		fmt.Printf("[WARN] There is a new version available. Would you like to update this program?")

		// ...
        // Ask user whether to update or not.
        // ...

        if shouldUpdate {
		    update(client, releaseProvider, conf)
        }
	}
}

func update(c http.HttpClientPlugin, p provider.UpdaterProvider, conf i18n.I18nConf) {
	err := caravela.Update(c, p, conf, "oalienista", "0.1.0")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("New version installed!")
	}
}
```

## Contributing
Please feel free to submit issues, fork the repository and send pull requests!

### Reporting issues
Bugs, feature requests, and development-related questions should be directed to our [GitHub issue tracker](https://github.com/aureliano/caravela/issues). If reporting a bug, please try and provide as much context as possible such as your operating system, Go version, and anything else that might be relevant to the bug. For feature requests, please explain what you're trying to do, and how the requested feature would help you do that.

Security related bugs can either be reported in the issue tracker.

### Submitting a patch
1. It's generally best to start by opening a new issue describing the bug or feature you're intending to fix. Even if you think it's relatively minor, it's helpful to know what people are working on. Mention in the initial issue that you are planning to work on that bug or feature so that it can be assigned to you.

2. Follow the normal process of [forking](https://help.github.com/articles/fork-a-repo) the project, and setup a new branch to work in. It's important that each group of changes be done in separate branches in order to ensure that a pull request only includes the commits related to that bug or feature.

3. Go makes it very simple to ensure properly formatted code, so always run `go fmt` on your code before committing it. You should also run `go vet` or `make code-lint` over your code. This will help you find common style issues within your code and will keep styling consistent within the project.

4. Any significant changes should almost always be accompanied by tests. The project already has good test coverage, so look at some of the existing tests if you're unsure how to go about it. [gocov](https://github.com/axw/gocov) and [gocov-html](https://github.com/matm/gocov-html) are invaluable tools for seeing which parts of your code aren't being exercised by your tests.

5. Please run:
 - `make test`
 - `make code-lint`

The `make test` command will run tests inside your code. This will help you spot places where code might be faulty before committing.

And the `make code-lint` command will check linting and styling over your code, keeping the project consistent formatting-wise.

Do your best to have [well-formed commit messages](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html) for each change. This provides consistency throughout the project, and ensures that commit messages are able to be formatted properly by various git tools.

Finally, push the commits to your fork and submit a [pull request](https://help.github.com/articles/creating-a-pull-request).

**NOTE**: Please do not use force-push on PRs in this repo, as it makes it more difficult for reviewers to see what has changed since the last code review.

### Other notes on code organization
Currently, all exposed features - check for updates and update - are defined in the main file at the base directory. Such features are delegated to providers according to where releases were published. That said, a Github provider must be a `provider/github.go` file as is `provider/gitlab.go` to Gitlab. So use that as your guide for where to put new providers.

### Maintainer's Guide
**Always try to maintain a clean, linear git history**. With very few exceptions, running `git log` should not show a bunch of branching and merging.

Never use the GitHub "merge" button, since it always creates a merge commit. Instead, check out the pull request locally ([these git aliases help](https://github.com/willnorris/dotfiles/blob/d640d010c23b1116bdb3d4dc12088ed26120d87d/git/.gitconfig#L13-L15)), then cherry-pick or rebase them onto master. If there are small cleanup commits, especially as a result of addressing code review comments, these should almost always be squashed down to a single commit. Don't bother squashing commits that really deserve to be separate though. If needed, feel free to amend additional small changes to the code or commit message that aren't worth going through code review for.

If you made any changes like squashing commits, rebasing onto master, etc, then GitHub won't recognize that this is the same commit in order to mark the pull request as "merged". So instead, amend the commit message to include a line "Fixes #0", referencing the pull request number. This would be in addition to any other "Fixes" lines for closing related issues. If you forget to do this, you can also leave a comment on the pull request. If you made any other changes, it's worth noting that as well.

## Versioning
This project follows the Semantic Versioning - [semver](https://semver.org/) - for tagging releases of the package.

 - The major version is incremented with any incompatible change.
 - The minor version is incremented with any backwards compatible change.
 - The patch version is incremented with any backwards compatible bug-fixes.

## License
This project is licensed under the terms of the MIT license found in the [LICENSE](./LICENSE) file.
