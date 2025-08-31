# oscar: the OpenSourceCorp Automation Runner

<!-- badges: start -->
![Github Actions](https://github.com/opensourcecorp/oscar/actions/workflows/main.yaml/badge.svg)
<!-- badges: end -->

`oscar` ("OpenSourceCorp Automation Runner") is a highly-opinionated, out-of-the-box task runner.
Originally designed for use exclusively within [OpenSourceCorp's CI/CD
subsystem](https://github.com/opensourcecorp/osc-infra/tree/main/cicd), it is perfectly usable
outside of OSC as well.

`oscar` is "highly-opinionated" in that `oscar` is designed to do each thing ***one single way***.
No choosing what linters to run or how to configure them, no picking which annual flavor of Python
or Nodejs packaging tool, etc. -- `oscar` is built to be ***the*** authoritative toolset for entire
teams and their codebases.

Under the hood, `oscar` uses the excellent [`mise`](https://mise.jdx.dev) quite heavily, and would
like to thank the author & contributors for making something like `oscar` possible without a lot of
wheel-reinvention.

## How to use

Before getting started, note that `oscar` has a few host-system runtime dependencies. Some of these
may someday be replaced natively in the future, but some are integral to how `oscar` works
internally.

* `bash` (version 4.4+)
* `git`

You run oscar by providing it a subcommand, such as `ci`. You can see the full available subcommand
list via `oscar --help`.

## Features

| Feature                      | `oscar` command | Details                            |
| :--------------------------- | :-------------- | :--------------------------------- |
| Continuous integration       | `oscar ci`      | [section](#continuous-integration) |
<!-- | Codebase & workstation setup | `oscar setup`   | [section]()                        | -->
<!-- | Delivery                     | `oscar deliver` | [section]()                        | -->
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

## Supported platforms

`oscar` is designed to run on Linux, and should work on macOS as well. Native Windows has not been
tested, and is unlikely to work. If you are on a Windows machine, you can run `oscar` in a WSL2
environment and it will work the same as on Linux.

## Development & Contributions

Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for details about developing `oscar`.

## Roadmap

* Doc updates
  * Functionality, rationale for decisions made, development history including switching to `mise`,
    etc.
* Workstation setup
  * Have `oscar` manage Makefiles, dotfiles, etc.
* CI additions
  * Terraform
  * protobuf
  * Rust?
* CD additions
  * Publish to ghcr
