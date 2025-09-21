ARG GO_VERSION
ARG MISE_VERSION

# These two proxy args are helpful if you're trying to build on a corporate network -- they do not
# impact the image if not, though.
ARG http_proxy
ARG https_proxy

####################################################################################################

FROM docker.io/library/golang:${GO_VERSION} AS builder

RUN apt-get update && apt-get install --no-install-recommends -y \
        bash \
        ca-certificates \
        make \
        && \
        rm -rf /var/lib/apt*

COPY . /go/app
WORKDIR /go/app

RUN make build

####################################################################################################

FROM docker.io/library/debian:13-slim AS ci

COPY --from=builder /go/app/build/oscar /oscar

RUN apt-get update && apt-get install --no-install-recommends -y \
        bash \
        ca-certificates \
        curl \
        git \
        gnupg2 \
        make \
        rename \
        && \
    rm -rf /var/lib/apt/*

COPY . /go/app
WORKDIR /go/app

# NOTE: when creating some shims, mise refers to itself assuming it is on the $PATH, so we need to
# symlink it out so it can do that
RUN ln -fs "${HOME}/.oscar/bin/mise" /usr/local/bin/mise && \
    bash ./scripts/test-bootstrap.sh setup && \
    /oscar ci

# ####################################################################################################

FROM docker.io/library/debian:13-slim AS final

ENV MISE_VERSION=${MISE_VERSION}

COPY --from=builder /go/app/build/oscar /oscar

# NOTE: Docker BuildKit will skip stages it doesn't see as dependencies, so to enforce the "ci"
# stage above to run, we need to force a dependency here
COPY --from=ci /go/app/LICENSE /LICENSE

RUN apt-get update && apt-get install --no-install-recommends -y \
        bash \
        ca-certificates \
        git \
        gnupg2 \
        && \
    rm -rf /var/lib/apt/*

# NOTE: when creating some shims, mise refers to itself assuming it is on the $PATH, so we need to
# symlink it out so it can do that
RUN ln -fs "${HOME}/.oscar/bin/mise" /usr/local/bin/mise

RUN groupadd --gid=1000 oscar && \
    useradd --uid=1000 --gid=1000 --create-home oscar && \
    mkdir -p /home/oscar/app
USER oscar
WORKDIR /home/oscar/app

# The location for the source code oscar will be run against
VOLUME /home/oscar/app
# oscar's home directory, for caching on the host
VOLUME /home/oscar/.oscar

# So e.g. GitHub can tie the image to its source repo
LABEL org.opencontainers.image.source https://github.com/opensourcecorp/oscar

ENTRYPOINT ["/oscar"]
