SHELL = /usr/bin/env bash -euo pipefail

BINNAME := oscar
BINPATH := ./cmd/$(BINNAME)

DOCKER ?= docker
OCI_REGISTRY ?= ghcr.io
OCI_REGISTRY_OWNER ?= opensourcecorp

SHELL = /usr/bin/env bash -euo pipefail

.PHONY: %

all: ci

ci: clean
	@go run ./cmd/oscar/main.go ci

# test is just an alias for ci
test: ci

ci-container:
	@docker build \
		--build-arg GO_VERSION="$$(awk '/^go/ { print $$2 }' go.mod)" \
		--build-arg CI=true \
		-f ./Containerfile -t $(BINNAME)-test:latest \
		.

build: clean
	@mkdir -p ./build/$$(go env GOOS)-$$(go env GOARCH)
	@go build -o ./build/$$(go env GOOS)-$$(go env GOARCH)/$(BINNAME) $(BINPATH)

xbuild: clean
	@for target in \
		darwin-amd64 \
		darwin-arm64 \
		linux-amd64 \
		linux-arm64 \
	; \
	do \
		GOOS=$$(echo "$${target}" | cut -d'-' -f1) ; \
		GOARCH=$$(echo "$${target}" | cut -d'-' -f2) ; \
		outdir=build/"$${GOOS}-$${GOARCH}" ; \
		mkdir -p "$${outdir}" ; \
		printf "Building for %s-%s into build/ ...\n" "$${GOOS}" "$${GOARCH}" ; \
		GOOS="$${GOOS}" GOARCH="$${GOARCH}" go build -o "$${outdir}"/$(BINNAME) $(BINPATH) ; \
	done

package: xbuild
	@mkdir -p dist
	@cd build || exit 1; \
	for built in * ; do \
		printf 'Packaging for %s into dist/ ...\n' "$${built}" ; \
		cd $${built} && tar -czf ../../dist/$(BINNAME)_$${built}.tar.gz * && cd - >/dev/null ; \
	done

clean:
	@rm -rf \
		/tmp/$(BINNAME)-tests \
		*cache* \
		.*cache* \
		*.log \
		build/ \
		dist/ \
		$(BINNAME) \
		./main

build-image: clean
	@$(DOCKER) build \
		--progress plain \
		--build-arg GO_VERSION="$$(awk '/^go/ { print $$2 }' go.mod)" \
		-f Containerfile \
		-t $(OCI_REGISTRY)/$(OCI_REGISTRY_OWNER)/$(BINNAME):latest \
		.
