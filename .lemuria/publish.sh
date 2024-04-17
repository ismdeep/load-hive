#!/usr/bin/env bash

set -e

# Get to workdir
cd "$(realpath "$(dirname "$(dirname "$(realpath "${BASH_SOURCE[0]}")")")")"

git config --global --add safe.directory "$(pwd)"

# 检查 Git 状态
git_status=$(git status --porcelain)

# 如果有未提交的变更，则输出提示信息
if [[ -n ${git_status} ]]; then
    echo "[WARN] There are changes that need to be committed:"
    echo "$git_status"
fi

# 获取当前 Git 提交信息
commit_date="$(TZ=Asia/Shanghai git log -1 --format=%cd --date=format:'%Y%m%d%H%M%S' | xargs echo)"
commit_id="$(git rev-parse --short=12 HEAD)"
version="$(cat VERSION)+${commit_date}-${commit_id}"

build-binary() {
  os="${1:?}" && \
  arch="${2:?}" && \
  mkdir -p "output/${commit_date}-${commit_id}/" && \
  source go.env.sh go1.22.1 && \
  GOOS=${os} GOARCH=${arch} GOPROXY=https://goproxy.cn,direct \
    go build \
      -o "output/${commit_date}-${commit_id}/load-hive-${os}-${arch}" \
      -trimpath \
      -ldflags '-s -w' . && \
  echo "OK: output/${commit_date}-${commit_id}/load-hive-${os}-${arch}"
}

# clean
rm -rf output/

# build
build-binary linux amd64
build-binary linux arm64

# publish
mkdir -p                                         "/data/public-data-work/data/load-hive/${version}/"
rsync -avz "output/${commit_date}-${commit_id}/" "/data/public-data-work/data/load-hive/${version}/"

# output
echo ""
echo "http://10.20.42.232/load-hive/${version}/"
echo ""
