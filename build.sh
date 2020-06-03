#!/bin/bash
# go build -ldflags "-X main.buildTime=$(date -u +"%Y.%m.%d.%H%M%S.UTC")"
GOLANG_VERSION="1.14.4"

docker run --rm -v "$PWD":/usr/local/go/src/hostingDeUpdateIP -w /usr/local/go/src/hostingDeUpdateIP golang:$GOLANG_VERSION /bin/bash -c 'go get ./... && go mod vendor && go build -ldflags "-X main.buildTime=$(date -u +"%Y.%m.%d.%H%M%S.UTC")"'
