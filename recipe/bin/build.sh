#!/bin/bash

set -euo pipefail

readonly BASEDIR="$(cd $(dirname $0)/.. && pwd)"
readonly TAG="${1?Provide a tag please}"

main() {

  build-binary downloader
  build-binary runner
  build-binary uploader

  build-packs-builder

  build-image downloader
  build-image runner
  build-image uploader

}

build-binary() {
  pushd "$BASEDIR/cmd"
    GOOS=linux go build -a -o "$BASEDIR"/image/downloader ${1}.go client.go
  popd
}

build-packs-builder() {
  pushd "$BASEDIR"/packs/cf/cmd/builder
    GOOS=linux CGO_ENABLED=0 go build -a -installsuffix static -o "$BASEDIR"/image/builder
  popd
}

build-image() {
  pushd "$BASEDIR"/image
    docker build --build-arg buildpacks="$(< "buildpacks.json")" -t "eirini/recipe:${TAG}" -f Dockerfile-${1} .
  popd
}

main
