# cloudcli

Terminal-based CLI interface for server and infrastructure management using supported APIs

## Install

Download binary for your OS/architecture from the [Releases](https://github.com/cloudwm/cloudcli/releases) page

```
sudo wget -O /usr/local/bin/cloudcli https://github.com/cloudwm/cloudcli/releases/download/v0.1.1/cloudcli-linux-amd64 &&\
sudo chmod +x /usr/local/bin/cloudcli
```

## Documentation

* [cloudcli](docs/cloudcli.md)	 - initial configuration and global options

#### Server management

* [cloudcli server create](docs/cloudcli_server_create.md)	 - Create a server
* [cloudcli server list](docs/cloudcli_server_list.md)	 - List servers
* [cloudcli server options](docs/cloudcli_server_options.md)	 - List server options
* [cloudcli server poweroff](docs/cloudcli_server_poweroff.md)	 - Power Off server/s
* [cloudcli server poweron](docs/cloudcli_server_poweron.md)	 - Power On server/s
* [cloudcli server terminate](docs/cloudcli_server_terminate.md)	 - Terminate server/s

#### Task queue management

* [cloudcli queue detail](docs/cloudcli_queue_detail.md)	 - Get details of tasks
* [cloudcli queue list](docs/cloudcli_queue_list.md)	 - List all tasks in queue
