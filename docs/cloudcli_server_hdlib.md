## cloudcli server hdlib

List cloneable hard-disks / Clone to hard-disk library 

### Synopsis

List cloneable hard-disks / Clone to hard-disk library 

```
cloudcli server hdlib [flags]
```

### Options

```
      --clone string        UUID Of hard-disk to clone
  -h, --help                help for hdlib
      --id string           Specific server UUID
      --image-name string   Name of image to clone the hard-disk as
      --name string         Server name or regular expression matching a single server
      --wait                Wait for command execution to finish only then exit cli.
```

### Options inherited from parent commands

```
      --api-clientid string   API Client ID
      --api-secret string     API Secret
      --config string         config file (default is $HOME/.cloudcli.yaml)
      --debug                 enable debug output to stderr
      --dryrun                enable dry run mode, does not perform actions
      --format string         output format, default format is a human readable summary
      --no-config             disable loading from config file
```

### SEE ALSO

* [cloudcli server](cloudcli_server.md)	 - Server management

###### Auto generated by spf13/cobra on 19-Oct-2024
