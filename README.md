# cloudcli

Terminal-based CLI interface for server and infrastructure management using supported APIs

## Download

Download the latest binary for your OS/architecture:

**Windows**: [64 bit](https://cloudcli.cloudwm.com/binaries/latest/cloudcli-windows-amd64.zip) | [32 bit](https://cloudcli.cloudwm.com/binaries/latest/cloudcli-windows-386.zip)

**Mac OS X**: [64 bit](https://cloudcli.cloudwm.com/binaries/latest/cloudcli-darwin-amd64.tar.gz) | [32 bit](https://cloudcli.cloudwm.com/binaries/latest/cloudcli-darwin-386.tar.gz)

**Linux**: [64 bit](https://cloudcli.cloudwm.com/binaries/latest/cloudcli-linux-amd64.tar.gz) | [32 bit](https://cloudcli.cloudwm.com/binaries/latest/cloudcli-linux-386.tar.gz)

## Install

Extract the downloaded archive and place the binary in your PATH.

Run `cloudcli init` to perform an interactive initialization.

## Usage

See the CLI help messages:

```
cloudcli --help
```

## Documentation

* [cloudcli](docs/cloudcli.md)	 - initial configuration and global options

#### Server management

* [cloudcli server create](docs/cloudcli_server_create.md)	 - Create a server
* [cloudcli server info](docs/cloudcli_server_info.md)	 - Get server overview/information
* [cloudcli server list](docs/cloudcli_server_list.md)	 - List servers
* [cloudcli server options](docs/cloudcli_server_options.md)	 - List server options
* [cloudcli server poweroff](docs/cloudcli_server_poweroff.md)	 - Power Off server/s
* [cloudcli server poweron](docs/cloudcli_server_poweron.md)	 - Power On server/s
* [cloudcli server terminate](docs/cloudcli_server_terminate.md)	 - Terminate server/s

#### Task queue management

* [cloudcli queue detail](docs/cloudcli_queue_detail.md)	 - Get details of tasks
* [cloudcli queue list](docs/cloudcli_queue_list.md)	 - List all tasks in queue
