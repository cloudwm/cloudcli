# cloudcli

Terminal-based CLI interface for server and infrastructure management using supported APIs

## Download

Download the latest binary for your OS/architecture:

**Windows**: [64 bit](https://cloudcli.cloudwm.com/binaries/latest/cloudcli-windows-amd64.zip) | [32 bit](https://cloudcli.cloudwm.com/binaries/latest/cloudcli-windows-386.zip)

**Mac OS X**: [64 bit](https://cloudcli.cloudwm.com/binaries/latest/cloudcli-darwin-amd64.zip)

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

* [cloudcli server attach](docs/cloudcli_server_attach.md)	 - Connect to the server via SSH
* [cloudcli server clone](docs/cloudcli_server_clone.md)	 - Clone a server
* [cloudcli server configure](docs/cloudcli_server_configure.md)	 - Change server configuration
* [cloudcli server create](docs/cloudcli_server_create.md)	 - Create a server
* [cloudcli server description](docs/cloudcli_server_description.md)	 - Get or set server description
* [cloudcli server disk](docs/cloudcli_server_disk.md)	 - List/manage server disks
* [cloudcli server history](docs/cloudcli_server_history.md)	 - List server actions history
* [cloudcli server info](docs/cloudcli_server_info.md)	 - Get server overview/information
* [cloudcli server list](docs/cloudcli_server_list.md)	 - List servers
* [cloudcli server network](docs/cloudcli_server_network.md)	 - List/manage server networks
* [cloudcli server options](docs/cloudcli_server_options.md)	 - List server options
* [cloudcli server passwordreset](docs/cloudcli_server_passwordreset.md)	 - Reset server/s password
* [cloudcli server poweroff](docs/cloudcli_server_poweroff.md)	 - Power Off server/s
* [cloudcli server poweron](docs/cloudcli_server_poweron.md)	 - Power On server/s
* [cloudcli server reboot](docs/cloudcli_server_reboot.md)	 - Reboot server/s
* [cloudcli server rename](docs/cloudcli_server_rename.md)	 - Rename server
* [cloudcli server reports](docs/cloudcli_server_reports.md)	 - Get server monthly usage reports
* [cloudcli server snapshot](docs/cloudcli_server_snapshot.md)	 - List/manage server snapshots
* [cloudcli server sshkey](docs/cloudcli_server_sshkey.md)	 - Add an SSH public key to the server authorized keys
* [cloudcli server statistics](docs/cloudcli_server_statistics.md)	 - Get server statistics
* [cloudcli server tags](docs/cloudcli_server_tags.md)	 - List/manage server tags
* [cloudcli server terminate](docs/cloudcli_server_terminate.md)	 - Terminate server/s

#### Task queue management

* [cloudcli queue detail](docs/cloudcli_queue_detail.md)	 - Get details of tasks
* [cloudcli queue list](docs/cloudcli_queue_list.md)	 - List all tasks in queue

#### Network management

* [cloudcli network create](docs/cloudcli_network_create.md)	 - Create a network
* [cloudcli network delete](docs/cloudcli_network_delete.md)	 - Delete a network (must delete all subnets first)
* [cloudcli network list](docs/cloudcli_network_list.md)	 - List networks
* [cloudcli network subnet_create](docs/cloudcli_network_subnet_create.md)	 - Create a network subnet
* [cloudcli network subnet_delete](docs/cloudcli_network_subnet_delete.md)	 - Delete a network subnet
* [cloudcli network subnet_edit](docs/cloudcli_network_subnet_edit.md)	 - Edit a network subnet, all values which are different then existing values will be updated
* [cloudcli network subnet_list](docs/cloudcli_network_subnet_list.md)	 - List network subnets
