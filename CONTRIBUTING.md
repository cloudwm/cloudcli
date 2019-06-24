# Contributing to cloudwm-cli

* Welcome to Cloudwm!
* Contributions of any kind are welcome.


## Building and running the CLI

We use Docker for consistent build environment. You are welcome to inspect `Dockerfile.build` and replicate on a local host.

Create the cloudcli configuration at /etc/cloudcli/.cloudcli.yaml - 
It's important to use this path as it's mounted inside the container.
See example-cloudcli.yaml file in the project root. 

Once the .cloudcli.yaml file is ready, you can build and start the build environment:

```
bin/build.sh start_build_environment
```

Compile and run the CLI from inside the build environment container:

```
bin/build.sh run --config /etc/cloudcli/.cloudcli.yaml server list
```

Build a binary so that it is available outside the container:

```
bin/build.sh build
```

Run the executable (From Linux):

```
./cloudcli
```


## Running the tests suite

* Compile the cloudcli binary and place in PATH
* Make sure `python` version 3.6 binary is available in PATH
* Make sure environment is "clean" - e.g. no default cloudcli config files / env vars.
* Set environment variables for testing:

```
TEST_API_SERVER=""
TEST_API_CLIENTID=""
TEST_API_SECRET=""
```

Verify the prerequisites

```
bin/test.sh verify
```

Run all tests:

```
bin/test.sh all
```
