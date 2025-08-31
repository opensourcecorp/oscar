SHELL = /usr/bin/env bash -euo pipefail

# Conditionally use mise if its available, otherwise expect the host to have any needed tools on $PATH
RUN =
MISE := $(shell command -v mise || command -v "$${HOME}/.oscar/bin/mise")
ifneq ($(MISE),)
RUN = $(MISE) exec --
endif

BINNAME := oscar
BINPATH := ./cmd/$(BINNAME)

DOCKER ?= docker
OCI_REGISTRY ?= ghcr.io
OCI_REGISTRY_OWNER ?= opensourcecorp

SHELL = /usr/bin/env bash -euo pipefail

.PHONY: %

all: ci

ci: clean
	@$(RUN) go run ./cmd/$(BINNAME)/main.go ci

# test is just an alias for ci
test: ci

ci-container:
	@$(DOCKER) build \
		--build-arg http_proxy="$${http_proxy}" \
		--build-arg https_proxy="$${https_proxy}" \
		--build-arg GO_VERSION="$$(awk '/^go/ { print $$2 }' go.mod)" \
		-f ./Containerfile \
		-t $(BINNAME)-test:latest \
		.

build: clean
	@mkdir -p ./build/$$($(RUN) go env GOOS)-$$($(RUN) go env GOARCH)
	@$(RUN) go build -o ./build/$(BINNAME) $(BINPATH)
	@printf 'built to %s\n' ./build/$(BINNAME)

xbuild: clean
	@bash ./scripts/xbuild.sh

package: xbuild

clean:
	@rm -rf \
		/tmp/$(BINNAME)-tests \
		./*cache* \
		./.*cache* \
		./*.log \
		./build/ \
		./dist/ \
		./$(BINNAME) \
		./main

image: clean
	@export BUILDKIT_PROGRESS=plain && \
	export GO_VERSION="$$(awk '/^go/ { print $$2 }' go.mod)" && \
	$(DOCKER) compose build

run-image:
	@$(DOCKER) compose run $(BINNAME)
