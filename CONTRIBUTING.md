# Development & Contribution

## Workstation setup

`oscar`'s development is the most consistent if you use the following:

* The root-level Makefile.
* [`mise`](https://mise.jdx.dev), via the root-level `mise.toml`.

While not strictly required, you will have a much better time if you leverage those. However, the
Makefile's targets are configured to also allow for working with most tooling natively, such as Go,
if they are available.

`mise` does not manage Docker Engine, however, so for targets like `make image` you will need to
ensure that you have a container runtime available (like the ones for Docker or Podman). The default
command the Makefile tries to run is `docker`, which you can override via the `DOCKER` Make
variable.

## Contribution philosophy

`oscar` welcomes contributions, so long as they adhere to a few key rules:

* As a reminder, `oscar`'s runtime behavior is ***intentionally designed to be rigid***. If there is
  a language or tool you would like to see added, then those contributions are welcome. Fundamental
  changes to how `oscar` intends to operate, e.g. allowing for user override of a linter's
  line-length checks, or which directory `oscar` builds binaries into, are not.

* `hEy tHeRe'S a LoT oF sHeLL cOdE iN hErE?!` -- Correct. "Shelling out" is an intentional design
  decision, and `oscar` calls out `bash` as a dependency.

## Adding support for a new Tool

To add a new Tool to `oscar`, it must:

* Define a non-exported struct, which embeds `taskutil.Tool`.
* That struct must implement the `taskutil.Tasker` interface.

See examples across the various `internal/tasks/tools/{lang}/*.go` files.
