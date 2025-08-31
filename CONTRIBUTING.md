# Development & Contribution

## Workstation setup

`oscar`'s development leans heavily into using the root-level Makefile.

Development of `oscar` is the most consistent if your workstation has [`mise`](https://mise.jdx.dev)
available on your `$PATH`. However, the Makefile's targets are configured to also allow for working
with most tooling natively, such as Go, if they are available.

`mise` does not manage Docker Engine, however, so for targets like `make image` you will need to
ensure that you have a container runtime available (like Docker or Podman, overridable via the
`DOCKER` Make variable).

## Contribution philosophy

`oscar` welcomes contributions, so long as they adhere to a few key rules:

* As a reminder, `oscar`'s runtime behavior is ***intentionally designed to be rigid***. If there is
  a language or tool you would like to see added, then those contributions are welcome. Fundamental
  changes to how `oscar` intends to operate, e.g. allowing for wholesale override of various
  linters' settings, are not.
