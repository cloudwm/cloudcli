export GOOS="${GOOS:-linux}"
export GOARCH="${GOARCH:-amd64}"

if [ "${1}" == "script" ]; then
    ( docker rm -f cloudwm-cli-build || true ) &&\
    docker build --build-arg GOOS=$GOOS --build-arg GOARCH=$GOARCH -t cloudwm-cli-build -f Dockerfile.build . &&\
    docker run -d --rm --name cloudwm-cli-build -v `pwd`:/go/src/github.com/cloudwm/cli cloudwm-cli-build tail -f /dev/null &&\
    sleep 1 &&\
    docker exec -it cloudwm-cli-build go build -o cloudcli main.go && sudo chown $USER ./cloudcli && sudo chmod +x ./cloudcli
    export PATH="`pwd`:${PATH}"
    [ "$?" != "0" ] && echo Failed build && exit 1
    if [ "${GOOS}" == "linux" ] && [ "${GOARCH}" == "amd64" ]; then
        echo "Running tests for linux/amd64"
        export DEBUG_OUTPUT_FILE=debug.log
        bash tests/test_all.sh
        RES="$?"
        cat debug.log
        exit "${RES}"
    else
        echo Skipping tests for $GOOS/$GOARCH
        exit 0
    fi
fi
