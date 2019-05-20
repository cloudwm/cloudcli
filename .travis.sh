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
    IMAGE="cloudwm/cloudcli:${TRAVIS_COMMIT}-${GOOS}-${GOARCH}"
    echo Pushing build environment to Docker image ${IMAGE}
    docker tag cloudwm-cli-build $IMAGE && docker push $IMAGE
    if [ "${GOOS}" == "linux" ] && [ "${GOARCH}" == "amd64" ]; then
        echo "Running tests for linux/amd64"
        # Debug output may contain sensitive details
        export DEBUG_OUTPUT_FILE=/dev/null
        bash tests/test_all.sh
        RES="$?"
        exit "${RES}"
    else
        echo Skipping tests for $GOOS/$GOARCH
        exit 0
    fi
fi
