# cloudcli

Terminal-based CLI interface for server and infrastructure management using supported APIs

## Install

Download binary for your OS/architecture from the [Releases](https://github.com/cloudwm/cloudcli/releases) page

```
sudo wget -O /usr/local/bin/cloudcli https://github.com/cloudwm/cloudcli/releases/download/v0.0.1/cloudcli-linux-amd64 &&\
sudo chmod +x /usr/local/bin/cloudcli
```

## Usage

Set server API host and credentials using one of the following options:

* **yaml configuration file**:
    * A configuration file can set CLI flags 
    * By default, a file is searched for at `$HOME/.cloudcli.yaml`
    * A different location can be specified using env var `CLOUDCLI_CONFIG=""` or `--config` flag
    * See [example-cloudcli.yaml](/example-cloudcli.yaml)

* **environment variables**:
    * Environment variables prefixed with `CLOUDCLI_` can set CLI flags
    * See [example-cloudcli.env](/example-cloudcli.env)

* **CLI flags**:
    * Server host and credentials can also be set using flags:
    * `cloudcli --api-server "" --api-clientid "" --api-secret ""`

**Important** Please keep your server and API credentials secure, 
it's recommended to use a configuration file with appropriate permissions and location.


## Commands

Following is an overview of main supported commands.

See the cloudcli command help messages for full reference.

```
cloudcli --help
```


### Server commands


#### `cloudcli server list`

List all the servers in the account


#### `cloudcli server options`

Download and show the available server options for your account.

One of the following flags must be provided:

```
  --billing      show billing resources
  --cpu          show cpu resources
  --datacenter   show datacenter resources
  --disk         show disk resources
  --image        show image resources
  --network      show network resources
  --ram          show ram resources
  --traffic      show traffic resources
```

If optional flag `--cache` is provided, the options will be downloaded to local file `cloudcli-server-options.json` and loaded from that file if it exists.


### Work in progress / Experiemental / Unstable commands

Following commands are work in progress / experiemental / unstable

To enable them set environment variable `CLOUDCLI_ENABLE_ALPHA=1`

##### `cloudcli server create`

Create server

**Work In Progress**
