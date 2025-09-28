# oscar: the OpenSourceCorp Automation Runner

<!-- badges: start -->
![Github Actions](https://github.com/opensourcecorp/oscar/actions/workflows/main.yaml/badge.svg)
<!-- badges: end -->

`oscar` ("OpenSourceCorp Automation Runner") is a highly-opinionated, out-of-the-box task runner
designed for use across OSC.

`oscar` is "highly-opinionated" in that it is designed to do each thing ***one single way***. No
choosing what linters to run or how to configure them, no picking which annual flavor of Python or
Nodejs packaging tool, no discrepancies in how to cut & deploy releases, etc. -- `oscar` is built to
be ***the*** authoritative toolset for entire teams and their codebases.

## Features

You run oscar by providing it a subcommand, such as `ci`. You can see the full available subcommand
list via `oscar --help`.

| Feature                | `oscar` command | Details                            |
| :--------------------- | :-------------- | :--------------------------------- |
| Continuous integration | `oscar ci`      | [section](#continuous-integration) |
| Delivery               | `oscar deliver` | [section](#delivery)               |
<!-- | Codebase & workstation setup | `oscar setup`   | [section]()                        | -->
<!-- | Deployment                   | `oscar deploy`  | [section]()                        | -->

### Continuous Integration

`oscar ci` runs a suite of continuous integration checks against your codebase, serving as something
of a linter aggregator. It provides these checks based on file discovery across your codebase, and
will only run checks based on what it finds. It also has behavior inspired by a tool named
[pre-commit](https://pre-commit.com/), including failing runs if any checks introduce Git diffs
during their runs.

Note again that these checks are ***highly opinionated*** -- if a particular linter supports
configuration, `oscar` configures it, but it uses ***its own built-in configuration***. For example,
if you try to change the line-length limit rule for `markdownlint-cli2` to be `120`, `oscar` will
ignore your request and run the check with its built-in limit of `100`. This behavior is
intentional, and serves to drive consistency across all manner of software that `oscar` could
possibly run against within a set of codebases.

However, this does not mean that someone is prevented from adding *additional* checks outside of
`oscar`'s purview -- it just means that you cannot override what `oscar` *does* control.

### Delivery

TODO

| Artifact types   | Targets          | `oscar.yaml` field               |
| :--------------- | :--------------- | :------------------------------- |
| Go binaries      | GitHub Releases  | `deliverables.go_github_release` |
| Container images | Any OCI registry | `deliverables.container_image`   |
<!-- | <empty cell>   | <second target for same artifact type> | -->
<!-- | <empty cell>   | <third target for same artifact type> | -->
<!-- | <second artifact type> | <first target for second artifact type> | -->

## Installation

`oscar` can be installed a few different ways:

* Downloadign a binary from a [GitHub Release](https://github.com/opensourcecorp/oscar/releases).

* Via `mise`, using a `github` or `go` backend:

      mise use "github:opensourcecorp/oscar@<version>"
      # or
      mise use "go:github.com/opensourcecorp/oscar/cmd/oscar@<version>"

* Via the Go toolchain:

      go install github.com/opensourcecorp/oscar/cmd/oscar@<version>

## Requirements

Before getting started, note that `oscar` has a few host-system runtime dependencies. Some of these
may someday be replaced natively in the future, but some are integral to how `oscar` works
internally.

* `bash` (version 4.4+)
* GNU `coreutils`
* `git`

In addition, some components of `oscar` may require additional host-system tools (e.g. a container
runtime like Docker for building & pushing container images).

If you are running on macOS, you should be able to install any missing tools via `brew install`-ing
the above by name -- but make sure your `$PATH` is pointing to the correct ones and not the default
BSD-equivalents.

## Supported platforms

`oscar` is designed to run on Linux, and should mostly work on macOS as well. Native Windows has not
been tested, and is unlikely to work. If you are on a Windows machine, you can run `oscar` in a WSL2
environment and it will work the same as on Linux.

## Development & Contributions

Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for details about developing `oscar`.

## Acknowledgements

Under the hood, `oscar` uses the excellent [`mise`](https://mise.jdx.dev) quite heavily, and would
like to thank the author & contributors for making something like `oscar` possible without a lot of
wheel-reinvention.

## Roadmap

* Add `oscar.yaml` generator
* Add check for changelog Markdown file that matches `oscar.yaml:version` (we should also use that
  file as the exact GH Release post contents)
* Workstation setup
  * Have `oscar` manage Makefiles, dotfiles, etc.
  * Also have it dump its own `mise.toml` for the user
  * `self-update` subcommand
* CI additions
  * Protobuf (especially since there's proto code in this repo now)
  * Terraform
* CD additions
  * Publishing to ghcr is confirmed to be working when run on `main` branch
