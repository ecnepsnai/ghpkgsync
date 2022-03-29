#!/bin/bash
set -e

VERSION=${1:?Version required}
REVISION=$(git rev-parse HEAD)
DATETIME=$(date --rfc-3339=seconds)

rm -f container/ghrpmsync server/ghrpmsync
cd server
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -buildmode=exe -o ghrpmsync
cd ../container
mv ../server/ghrpmsync .
podman build \
    --squash \
    --no-cache \
    --format docker \
    --label "org.opencontainers.image.created=${DATETIME}" \
    --label "org.opencontainers.image.version=${VERSION}" \
    --label "org.opencontainers.image.revision=${REVISION}" \
    -t ghcr.io/ecnepsnai/ghrpmsync:latest \
    -t ghcr.io/ecnepsnai/ghrpmsync:${VERSION} \
    .
rm -f ghrpmsync
