#!/usr/bin/env bash
set -euo pipefail
shopt -s globstar

_setup() {
  cp ./internal/ci/configfiles/pyproject.toml .
  mkdir -p ./src
  cp -r ./testdata/python/src ./src/test_package
  rename 's/\.test//g' -- ./src/**
}

_teardown() {
  rm -rf \
    ./pyproject.toml \
    ./src
}

cmd="${1:-}"

case "${cmd}" in
  setup)
    _setup
  ;;
  teardown)
    _teardown
  ;;
  *)
    printf '"%s" is not a recognized command\n' "${cmd}" >&2
    exit 1
  ;;
esac
