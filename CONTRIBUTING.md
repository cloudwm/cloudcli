# Contributing to cloudwm-cli

* Welcome to Cloudwm!
* Contributions of any kind are welcome.


## Building and running the CLI

We use Docker for consistent build environment. You are welcome to inspect `Dockerfile.build` and replicate on a local host.

Build and run the Docker image which contains the build environment:

```
( docker rm -f cloudwm-cli-build || true ) &&\
docker build -t cloudwm-cli-build -f Dockerfile.build . &&\
docker run -d --rm --name cloudwm-cli-build -v `pwd`:/go/src/github.com/cloudwm/cli cloudwm-cli-build tail -f /dev/null
```

Compile and run the CLI:

```
docker exec -it cloudwm-cli-build go run main.go
```

Use [cobra](https://github.com/spf13/cobra/blob/master/cobra/README.md) to add commands:

```
docker exec -it cloudwm-cli-build cobra add --help
```

Build a Linux binary and set executable:

```
docker exec -it cloudwm-cli-build go build main.go && sudo chown $USER ./main && sudo chmod +x ./main
```

Run the executable (From Linux):

```
./main
```

(Optional) For fast development iterations, define bash aliases:

```
alias cloudcli="docker exec -it cloudwm-cli-build go run main.go"
alias cloudcli-cobra="docker exec -it cloudwm-cli-build cobra add --help"
alias cloudcli-build="docker exec -it cloudwm-cli-build go build main.go && sudo chown $USER ./main && sudo chmod +x ./main"
```

## Troubleshooting

If you encounter permission problems, try running the following from project path in local host:

```
sudo chown -R $USER .
```


## Cross Platform Building

Build cross platform binaries for Windows, Mac and Linux:

```
for GOOS in darwin linux windows; do
  for GOARCH in 386 amd64; do
    echo "Building $GOOS-$GOARCH"
    export GOOS=$GOOS
    export GOARCH=$GOARCH
    docker build --build-arg GOOS=$GOOS --build-arg GOARCH=amd64 -t cloudwm-cli-build-$GOOS-$GOARCH -f Dockerfile.build . &&\
    docker run -it -v `pwd`:/go/src/github.com/cloudwm/cli cloudwm-cli-build-$GOOS-$GOARCH go build -o main-$GOOS-$GOARCH main.go
  done
done
```
