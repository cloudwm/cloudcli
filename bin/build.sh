#!/usr/bin/env bash

export BUILD_CONTAINER_NAME="${BUILD_CONTAINER_NAME:-cloudwm-cli-build}"
export GOOS="${GOOS:-linux}"
export GOARCH="${GOARCH:-amd64}"
export CLOUDCLI_ETC_PATH="${CLOUDCLI_ETC_PATH:-/etc/cloudcli}"
if [ "${CLOUDCLI_PROJECT_PATH}" == "" ]; then
    export CLOUDCLI_PROJECT_PATH="$(pwd)"
else
    export CLOUDCLI_PROJECT_PATH="${CLOUDCLI_PROJECT_PATH}"
fi

start_build_environment() {
    ( docker rm -f "${BUILD_CONTAINER_NAME}" || true ) &&\
    docker build --build-arg GOOS="${GOOS}" --build-arg GOARCH="${GOARCH}" \
                 -t "${BUILD_CONTAINER_NAME}" \
                 -f "${CLOUDCLI_PROJECT_PATH}/Dockerfile.build" "${CLOUDCLI_PROJECT_PATH}" &&\
    docker run -d --rm --name "${BUILD_CONTAINER_NAME}" -v "${CLOUDCLI_PROJECT_PATH}:/go/src/github.com/cloudwm/cli" \
               -v "${CLOUDCLI_ETC_PATH}:${CLOUDCLI_ETC_PATH}" --network host \
               "${BUILD_CONTAINER_NAME}" tail -f /dev/null &&\
    docker exec -it "${BUILD_CONTAINER_NAME}" dep ensure
}

run() {
    docker exec -it "${BUILD_CONTAINER_NAME}" go run main.go "$@"
}

build() {
    docker exec -it "${BUILD_CONTAINER_NAME}" go build -o cloudcli main.go &&\
    sudo chown $USER ./cloudcli &&\
    sudo chmod +x ./cloudcli
}

CMD="${1:-all}"

if [ "${CMD}" == "all" ]; then
    start_build_environment &&\
    run --config "${CLOUDCLI_ETC_PATH}/.cloudcli.yaml" &&\
    build &&\
    ./cloudcli --config "${CLOUDCLI_ETC_PATH}/.cloudcli.yaml" &&\
    ls -lah ./cloudcli
elif [ "${CMD}" == "start_build_environment" ]; then
    start_build_environment
elif [ "${CMD}" == "run" ]; then
    run "${@:2}"
elif [ "${CMD}" == "build" ]; then
    build
else
    echo "Invalid ARG: ${CMD}"
    false
fi && echo "Great Success!"
