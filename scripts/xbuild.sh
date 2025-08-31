#!/usr/bin/env bash
set -euo pipefail

targets=$(go tool dist list | grep -E 'linux|darwin' | grep -E 'amd64|arm64')
printf 'Will build for:\n'
while read -r line ; do
  printf '\t%s\n' "${line}"
done <<< "${targets}"

for target in ${targets} ; do
  GOOS=$(echo "${target}" | cut -d'/' -f1)
  GOARCH=$(echo "${target}" | cut -d'/' -f2)
  export GOOS GOARCH

  mkdir -p ./build
  out="./build/oscar-${GOOS}-${GOARCH}"
  printf "Building to %s\n" "${out}"
  go build -o "${out}" ./cmd/oscar
  chmod +x "${out}"
done
