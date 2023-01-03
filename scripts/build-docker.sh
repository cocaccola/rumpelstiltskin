#!/usr/bin/env bash

set -eu -o pipefail

source scripts/version.sh

docker build -t "rumpelstiltskin:$version" .
