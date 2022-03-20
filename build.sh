#!/bin/bash
set -e

rm -f container/ghrpmsync server/ghrpmsync
cd server
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o ghrpmsync
cd ../container
mv ../server/ghrpmsync .
podman build \
    --squash \
    --no-cache \
    --format docker \
    --label "org.opencontainers.image.created=$(date --rfc-3339=seconds)" \
    -t ghcr.io/ecnepsnai/ghrpmsync:latest \
    .
rm -f ghrpmsync
