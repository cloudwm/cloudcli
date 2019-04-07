# Contributing to cloudwm-cli

* Welcome to Cloudwm!
* Contributions of any kind are welcome.


## Building and running the CLI

We use Docker for consistent build environment. You are welcome to inspect `Dockerfile.build` and replicate on a local host.

Set env vars for desired OS (doesn't have to match your own OS)

```
export GOOS=linux
export GOARCH=amd64
```

Build and run the Docker image which contains the build environment:

```
( docker rm -f cloudwm-cli-build || true ) &&\
docker build --build-arg GOOS=$GOOS --build-arg GOARCH=$GOARCH -t cloudwm-cli-build -f Dockerfile.build . &&\
docker run -d --rm --name cloudwm-cli-build -v `pwd`:/go/src/github.com/cloudwm/cli \
           -v /etc/cloudcli:/etc/cloudcli \
           cloudwm-cli-build tail -f /dev/null
```

Compile and run the CLI:

```
docker exec -it cloudwm-cli-build go run main.go
```

(Optional) Enable alpha commands and set a configuration file:

```
docker exec -e CLOUDCLI_CONFIG=/etc/cloudcli/.my-config.yaml -e CLOUDCLI_ENABLE_ALPHA=1 -it cloudwm-cli-build go run main.go
```

(Optional) For fast development iterations, define bash aliases:

```
alias cloudcli="docker exec -it cloudwm-cli-build go run main.go"
alias cloudcli-build="docker exec -it cloudwm-cli-build go build -o cloudcli main.go && sudo chown $USER ./cloudcli && sudo chmod +x ./cloudcli"
```


Build a binary and set executable:

```
docker exec -it cloudwm-cli-build go build -o cloudcli main.go && sudo chown $USER ./cloudcli && sudo chmod +x ./cloudcli
```

Run the executable (From Linux):

```
./cloudcli
```

## Troubleshooting

If you encounter permission problems, try running the following from project path in local host:

```
sudo chown -R $USER .
```

Environment variables are not passed to the container by default. To test environment variables - you will have to modify the aliases or modify the exec script.


## Cross Platform Building

Build cross platform binaries for Windows, Mac and Linux:

```
for GOOS in darwin linux windows; do
  for GOARCH in 386 amd64; do
    echo "Building $GOOS-$GOARCH"
    export GOOS=$GOOS
    export GOARCH=$GOARCH
    docker build --build-arg GOOS=$GOOS --build-arg GOARCH=amd64 -t cloudwm-cli-build-$GOOS-$GOARCH -f Dockerfile.build . &&\
    docker run -it -v `pwd`:/go/src/github.com/cloudwm/cli cloudwm-cli-build-$GOOS-$GOARCH go build -o cloudcli-$GOOS-$GOARCH main.go
  done
done
```


## Running the tests suite

Compile the cloudcli binary and place in PATH

Make sure environment is "clean" - e.g. no default cloudcli config files / env vars.

Set environment variables for the API server and account to use for testing:

```
TEST_API_SERVER=""
TEST_API_CLIENTID=""
TEST_API_SECRET=""
```

Run the tests suite:

```
tests/test_all.sh
```

## Build Environments History

Each build does fresh pulls of system dependencies and go modules.

To revert or debug previous build environments, travis-ci pushed the build image to Docker Hub.

The image tag is `cloudwm/cloudcli:{GIT_COMMIT_SHA}-{GOOS}-{GOARCH}` 
