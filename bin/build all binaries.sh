#!/usr/bin/env bash

( [ "${GIT_REPO_URL}" == "" ] || [ "${GIT_BRANCH}" == "" ]  || [ "${BUILD_ENV_DOCKER_IMAGE_BASE_NAME}" == "" ] || [ "${BUILD_ENV_DOCKER_IMAGE_TAG}" == "" ] ) \
    && echo missing required env vars && exit 1

if [ "${BUILD_ENV_DOCKER_IMAGE_TAG}" == "master" ]; then
  export BUILD_ENV_DOCKER_IMAGE_TAG="latest"
fi

build() {
  local GOOS="${1}"
  local GOARCH="${2}"
  local EXT="${3}"
  local DOCKER_IMAGE="${BUILD_ENV_DOCKER_IMAGE_BASE_NAME}${GOOS}-${GOARCH}:${BUILD_ENV_DOCKER_IMAGE_TAG}"
  if docker pull "${DOCKER_IMAGE}"; then
    docker build --cache-from "${DOCKER_IMAGE}" --build-arg GOOS="${GOOS}" --build-arg GOARCH="${GOARCH}" \
                 -t "${DOCKER_IMAGE}" -f ./Dockerfile.build .
  else
    docker build -t "${DOCKER_IMAGE}" -f ./Dockerfile.build .
  fi &&\
  docker push "${DOCKER_IMAGE}" &&\
  docker run --rm -v "`pwd`:/go/src/github.com/cloudwm/cli" "${DOCKER_IMAGE}" go build -o cloudcli-${GOOS}-${GOARCH}${EXT} main.go
}

build darwin 386 "" &&\
build darwin amd64 "" &&\
build linux 386 "" &&\
build linux amd64 "" &&\
build windows 386 ".exe" &&\
build windows amd64 ".exe" &&\
echo Great Success!
