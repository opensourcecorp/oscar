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
export IMAGE_REGISTRY ?= ghcr.io
export IMAGE_REGISTRY_OWNER ?= opensourcecorp
export IMAGE_NAME ?= $(BINNAME)
export IMAGE_TAG ?= latest
export IMAGE_URI ?= $(IMAGE_REGISTRY)/$(IMAGE_REGISTRY_OWNER)/$(IMAGE_NAME):$(IMAGE_TAG)

SHELL = /usr/bin/env bash -euo pipefail

all: ci

FORCE:

ci: clean
	@$(RUN) go run ./cmd/$(BINNAME)/main.go ci

deliver:
	@$(RUN) go run ./cmd/$(BINNAME)/main.go deliver

# test is just an alias for ci
test: ci

# NOTE: oscar builds itself IRL, but having a target here makes it easier to have the Containerfile
# have a stage-copiable output
build: FORCE
	@$(RUN) go build -o ./build/oscar ./cmd/oscar

clean: FORCE
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
	$(RUN) $(DOCKER) compose build

run-image: FORCE
	@$(RUN) $(DOCKER) compose run $(BINNAME)

generate: FORCE
	@cd ./proto && $(RUN) buf generate
