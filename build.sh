#!/bin/bash
set -e

VERSION=${1:?Version required}
REVISION=$(git rev-parse HEAD)
DATETIME=$(date --rfc-3339=seconds -u)

rm -f container/ghpkgsync server/ghpkgsync
cd server
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOAMD64=v2 go build -ldflags="-s -w -X 'main.Version=${VERSION}'" -trimpath -buildmode=exe
cd ../container
mv ../server/ghpkgsync .
podman build \
    --squash \
    --no-cache \
    --format docker \
    --label "org.opencontainers.image.created=${DATETIME}" \
    --label "org.opencontainers.image.version=${VERSION}" \
    --label "org.opencontainers.image.revision=${REVISION}" \
    -t ghcr.io/ecnepsnai/ghpkgsync:latest \
    -t ghcr.io/ecnepsnai/ghpkgsync:${VERSION} \
    .
rm -f ghpkgsync
