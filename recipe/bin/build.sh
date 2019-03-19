#!/bin/bash

set -euo pipefail

readonly BASEDIR="$(cd $(dirname $0)/.. && pwd)"
readonly TAG="${1?Provide a tag please}"

main() {

  build-binary downloader
  build-binary executor
  build-binary uploader

  build-packs-builder

  build-image downloader
  build-image executor
  build-image uploader
}

build-binary() {
  pushd "${BASEDIR}/cmd/${1}"
    GOOS=linux go build -a -o "$BASEDIR"/image/${1} *.go
  popd
}

build-packs-builder() {
  pushd "$BASEDIR"/packs/cf/cmd/builder
    GOOS=linux CGO_ENABLED=0 go build -a -installsuffix static -o "$BASEDIR"/image/builder
  popd
}

# TODO: cleanup build image and build args
build-image() {
  pushd "$BASEDIR"/image
    docker build --build-arg buildpacks="$(< "buildpacks.json")" -t "eirini/recipe-${1}:${TAG}" -f Dockerfile-${1} .
  popd
}

main
