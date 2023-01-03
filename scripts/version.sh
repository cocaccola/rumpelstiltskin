#!/usr/bin/env bash

set -eu -o pipefail

source VERSION

current_commit=$(git rev-parse --short HEAD)
current_branch=$(git name-rev --name-only HEAD)
current_branch=${current_branch##refs/heads/}

if [[ -z $VERSION || current_branch != "main" || current_branch != "master" ]]; then
    version="experimental-$(date +%F-%H-%M)-$current_commit"
else
    version=$VERSION
fi

export version
