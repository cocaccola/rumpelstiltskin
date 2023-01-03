#!/usr/bin/env bash

set -eu -o pipefail -x

source scripts/version.sh

arch=$1
if [[ $arch != "arm64" && $arch != "amd64" ]]; then
    echo "invalid arch" >&2
    exit 2
fi


GOOS=linux GOARCH=$arch go build \
    -ldflags="-X main.version=$version -extldflags=-static -s -w" \
    -o rumpelstiltskin \
    main.go
