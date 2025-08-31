# This Containerfile isn't mean for creating artifacts etc., it's just a way to perform portable,
# local CI checks in case there are workstation-specific issues a developer faces.
FROM debian:13-slim AS builder

ARG MISE_VERSION="v2025.8.21"

# These two proxy args are helpful if you're trying to build on a corporate network -- they do not
# impact the image if not, though.
ARG http_proxy
ARG https_proxy

RUN apt-get update && apt-get install -y \
        bash \
        ca-certificates \
        curl \
        git \
        make \
        rename \
        tar \
        xz-utils

COPY . /go/app
WORKDIR /go/app

RUN curl -fsSL -o /usr/local/bin/mise "https://github.com/jdx/mise/releases/download/${MISE_VERSION}/mise-${MISE_VERSION}-linux-x64" && \
    chmod +x /usr/local/bin/mise && \
    mise trust --all && \
    mise install

# NOTE: we want to run CI twice -- once to make sure it works, and another to make sure it's
# *faster* because of the version-checking
RUN bash ./scripts/test-handler.sh setup && \
    bash -c 'echo $(date +%s) > /tmp/starttime' && \
    make ci && \
    bash -c 'echo $(date +%s) > /tmp/endtime' && \
    bash -c 'echo $(( $(cat /tmp/endtime) - $(cat /tmp/starttime) )) > /tmp/duration_1'
RUN bash -c 'echo $(date +%s) > /tmp/starttime' && \
    make ci && \
    bash -c 'echo $(date +%s) > /tmp/endtime' && \
    bash -c 'echo $(( $(cat /tmp/endtime) - $(cat /tmp/starttime) )) > /tmp/duration_2'
RUN bash -c 'if [[ ! "$(cat /tmp/duration_2)" -lt "$(cat /tmp/duration_1)" ]] ; then echo "Second CI run should have taken less time than the first" && exit 1 ; fi'

# Build last so the final image has access to copy the binary
RUN make build

###############################################################################

FROM docker.io/library/debian:13-slim

COPY --from=builder /go/app/build/oscar /oscar

RUN apt-get update && apt-get install -y \
        bash \
        ca-certificates \
        curl \
        git \
        make \
        rename \
        tar \
        xz-utils

RUN groupadd --gid=1000 oscar && \
    useradd --uid=1000 --gid=1000 --create-home oscar && \
    mkdir -p /home/oscar/app
USER oscar
WORKDIR /home/oscar/app

RUN curl -fsSL -o /usr/local/bin/mise "https://github.com/jdx/mise/releases/download/${MISE_VERSION}/mise-${MISE_VERSION}-linux-x64" && \
    chmod +x /usr/local/bin/mise && \
    mise trust --all

VOLUME /home/oscar/app
VOLUME /home/oscar/.oscar

ENTRYPOINT ["/oscar"]
