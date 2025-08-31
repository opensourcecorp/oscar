# oscar

<!-- badges: start -->
![Github Actions](https://github.com/opensourcecorp/oscar/actions/workflows/main.yaml/badge.svg)
<!-- badges: end -->

`oscar` ("OpenSourceCorp Automation Runner") is a CI/CD task runner designed for use within
[OpenSourceCorp's CI/CD subsystem](https://github.com/opensourcecorp/osc-infra/tree/main/cicd).

`oscar` does not orchestrate CI/CD tasks -- that's the subsystem's job. Rather, it is a set of
utilities that are designed to be easily ported between any CI/CD platform. `oscar` comprises most
of the CI/CD logic that your platform would normally run as steps in that process -- test, build,
push, deploy, etc. In this way, you can think of `oscar` being much like a Jenkins shared library.

## How to use

Before getting started, note that `oscar` has a few host-system runtime dependencies. Some of these
may someday be replaced natively in the future, but some are integral to how `oscar` works
internally.

* [`mise`](https://mise.jdx.dev) (though this will be installed for you if it is missing)
* `bash` (version 4.4+)
* `git`

TODO:

* Check if `mise` is installed, and install it if not
* Run as image from ghcr
* etc.
