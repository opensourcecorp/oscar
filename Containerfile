# This Containerfile isn't mean for creating artifacts etc., it's just a way to perform portable,
# local CI checks in case there are workstation-specific issues a developer faces.
ARG GO_VERSION=UNSET
FROM docker.io/library/golang:${GO_VERSION}

ARG CI

RUN apt-get update && apt-get install -y \
    ca-certificates \
    make

COPY . /go/app
WORKDIR /go/app

RUN if [ -n "${CI}" ] ; then make ci ; fi
