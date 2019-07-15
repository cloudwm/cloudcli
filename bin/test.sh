#!/usr/bin/env bash

if [ "${CLOUDCLI_PROJECT_PATH}" == "" ]; then
    export CLOUDCLI_PROJECT_PATH="$(pwd)"
else
    export CLOUDCLI_PROJECT_PATH="${CLOUDCLI_PROJECT_PATH}"
fi

verify() {
    local python_path="$(which python)"
    local python_version="$(python --version)"
    local cloudcli_path="$(which cloudcli)"
    echo python_path=$python_path
    echo python_version=$python_version
    echo cloudcli_path=$cloudcli_path
    echo test_api_server=$TEST_API_SERVER
    if [ "${python_path}" == "" ]; then
        echo Missing Python interpreter, make sure Python 3.6 is available as python in your PATH
    elif ! echo $python_version | grep 3.6 >/dev/null 2>&1; then
        echo Invalid Python interpreter, make sure Python 3.6 is available as python in your PATH
        false
    elif [ "${cloudcli_path}" == "" ]; then
        echo cloudcli binary must be in your PATH
        false
    elif [ -z "${TEST_API_SERVER}" ]; then
        echo Missing required environment variable: TEST_API_SERVER
        false
    elif [ -z "${TEST_API_CLIENTID}" ]; then
        echo Missing required environment variable: TEST_API_CLIENTID
        false
    elif [ -z "${TEST_API_SECRET}" ]; then
        echo Missing required environment variable: TEST_API_SECRET
        false
    fi
}

CMD="$1"

if [ "${CMD}" == "verify" ]; then
    verify
elif [ "${CMD}" == "all" ] || [ "${CMD}" == "" ]; then
    verify && tests/test_all.sh
else
    echo invalid argument: $CMD
    false
fi && echo Great Success!
