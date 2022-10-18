## cloudcli network subnet_edit

Edit a network subnet, all values which are different then existing values will be updated

### Synopsis

Edit a network subnet, all values which are different then existing values will be updated

```
cloudcli network subnet_edit [flags]
```

### Options

```
      --datacenter string          (required) edit a subnet in this datacenter
      --dns1 string                (optional) dns1 (e.g. 1.2.3.4)
      --dns2 string                (optional) dns2 (e.g. 2.3.4.5)
      --gateway string             (optional) gateway (e.g. 172.16.0.100)
  -h, --help                       help for subnet_edit
      --subnetBit string           (required) subnetBit (e.g. 23)
      --subnetDescription string   (optional) subnet description
      --subnetId string            (required) the subnet id to edit
      --subnetIp string            (required) subnetIP (e.g. 172.16.0.0)
      --vlanId string              (required) the vlan id the subnet belongs to
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

* [cloudcli network](cloudcli_network.md)	 - Network management

###### Auto generated by spf13/cobra on 22-Sep-2022