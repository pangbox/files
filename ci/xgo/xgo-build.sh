#!/bin/sh
set -e -u -o pipefail -x

XGO_DOCKER_IMAGE=${1:-ghcr.io/pangbox/pangfiles/xgo:latest}

mkdir -p bin
go install github.com/crazy-max/xgo@v0.32.0
xgo --docker-image="${XGO_DOCKER_IMAGE}" --targets=darwin/amd64 --pkg=cmd/pang --out bin/pang .
xgo --docker-image="${XGO_DOCKER_IMAGE}" --targets=darwin/arm64 --pkg=cmd/pang --out bin/pang .
