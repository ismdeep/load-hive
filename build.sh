#!/usr/bin/env bash

set -e

# Get to workdir
cd "$(realpath "$(dirname "$(realpath "${BASH_SOURCE[0]}")")")"

VERSION="$(cat "VERSION")"
COMMIT_ID="$(git rev-parse --short=12 HEAD)"
COMMIT_DATE="$(TZ=Asia/Shanghai git log -1 --format=%cd --date=format:'%Y%m%d%H%M%S' | xargs echo)"

output_folder="output/load-hive-${VERSION}-${COMMIT_DATE}-${COMMIT_ID}"

build-binary() {
  os="${1:?}" && \
  arch="${2:?}" && \
  GOOS="${os}" GOARCH="${arch}" go build -o "${output_folder}/load-hive-${os}-${arch}" -mod vendor -trimpath -ldflags '-s -w' . && \
  echo "[OK] ${output_folder}/load-hive-${os}-${arch}"
}

rm -rf "${output_folder:?}/"
build-binary darwin amd64
build-binary darwin arm64
build-binary linux  amd64
build-binary linux  arm64
