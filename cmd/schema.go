package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var schemaJson = `{
  "schema_version": [0, 0, 2],
  "commands": [
    {
      "use": "server",
      "short": "Server management",
      "commands": [
        {
          "use": "list",
          "short": "List servers",
          "run": {
            "cmd": "getList",
            "path": "/service/servers",
            "fields": [
              {
                "name": "id"
              },
              {
                "name": "name"
              },
              {
                "name": "datacenter"
              },
              {
                "name": "power"
              }
            ]
          }
        },
        {
          "use": "options",
          "short": "List server options",
          "run": {
            "cmd": "getListOfLists",
            "path": "/service/server",
            "lists": [
              {
                "name": "datacenter",
                "key": "datacenters",
                "type": "map",
                "fields": [
                  {
                    "name": "datacenter",
                    "from": "key"
                  },
                  {
                    "name": "datacenter_description",
                    "from": "value",
                    "parsers": [
                      {
                        "parser": "split_value_remove_first",
                        "split_string": ":"
                      }
                    ]
                  }
                ]
              },
              {
                "name": "cpu",
                "key": "cpu",
                "type": "list",
                "fields": [
                  {"name":  "cpu_type", "from":  "value"}
                ]
              },
              {
                "name": "ram",
                "key": "ram",
                "type": "list",
                "fields": [
                  {"name":  "ram_size_gb", "from":  "value"}
                ]
              },
              {
                "name": "disk",
                "key": "disk",
                "type": "list",
                "fields": [
                  {"name":  "disk_size_gb", "from":  "value"}
                ]
              },
              {
                "name": "image",
                "key": "diskImages",
                "type": "mapOfLists",
                "fields": [
                  {"name":  "image_id", "from":  "id"},
                  {"name":  "image_name", "from":  "description"},
                  {"name":  "image_size_gb", "from":  "sizeGB"},
                  {"name":  "image_usage_info", "from":  "usageInfo", "long": true},
                  {"name":  "image_guest_description", "from":  "guestDescription", "long": true},
                  {"name":  "image_text_one", "from":  "freeTextOne", "long": true},
                  {"name":  "image_text_two", "from":  "freeTextTwo", "long": true}
                ]
              },
              {
                "name": "traffic",
                "key": "traffic",
                "type": "mapOfLists",
                "fields": [
                  {"name":  "datacenter", "from":  "key"},
                  {"name":  "traffic", "from":  "id"},
                  {"name":  "traffic_info", "from":  "info"}
                ]
              },
              {
                "name": "network",
                "key": "networks",
                "type": "mapOfLists",
                "fields": [
                  {"name":  "datacenter", "from":  "key"},
                  {"name":  "network", "from":  "name"},
                  {
                    "name":  "network_ips", "from":  "ips",
                    "parsers": [
                      {"parser": "network_ips", "only_for_humans":  true}
                    ]
                  }
                ]
              },
              {
                "name": "billing",
                "key": "billing",
                "type": "list",
                "fields": [
                  {"name":  "billing_plan", "from":  "value"}
                ]
              }
            ]
          }
        },
        {
          "alpha": true,
          "use": "create",
          "short": "Create a server",
          "flags": [
            {
              "name": "name",
              "usage": "Server name (a-zA-Z0-9()_-). (must be at least 4 characters long, mandatory)",
              "required": true
            },
            {
              "name": "datacenter",
              "usage": "Server datacenter (EU, US-NY2, AS.. see --list-options). (mandatory)",
              "required": true
            },
            {
              "name": "image",
              "usage": "Server image name or image ID (see --list-options). (mandatory)",
              "required": true
            }
          ],
          "run": {
            "cmd": "post",
            "path": "/service/server",
            "fields": [
              {
                "name": "name",
                "flag": "name"
              },
              {
                "name": "datacenter",
                "flag": "datacenter"
              },
              {
                "name": "image",
                "flag": "image"
              }
            ]
          }
        }
      ]
    }
  ]
}`

type SchemaCommandFieldParser struct {
	Parser string `json:"parser"`
	SplitString string `json:"split_string"`
	OnlyForHumans bool `json:"only_for_humans"`
}

type SchemaCommandField struct {
	Name string `json:"name"`
	Flag string `json:"flag"`
	From string `json:"from"`
	Parsers []SchemaCommandFieldParser `json:"parsers"`
	Long bool `json:"long"`
}

type SchemaCommandFlag struct {
	Name string `json:"name"`
	Usage string `json:"usage"`
	Required bool `json:"required"`
}

type SchemaCommandList struct {
	Name string `json:"name"`
	Key string `json:"key"`
	Type string `json:"type"`
	Fields []SchemaCommandField `json:"fields"`

}

type SchemaCommandRun struct {
	Cmd string `json:"cmd"`
	Path string `json:"path"`
	Fields []SchemaCommandField `json:"fields"`
	Lists []SchemaCommandList `json:"lists"`
}

type SchemaCommand struct {
	Alpha bool `json:"alpha"`
	Use string `json:"use"`
	Short string `json:"short"`
	Run SchemaCommandRun `json:"run"`
	Flags []SchemaCommandFlag `json:"flags"`
	Commands []SchemaCommand `json:"commands"`
}

type Schema struct {
	SchemaVersion [3]int `json:"schema_version"`
	Commands []SchemaCommand `json:"commands"`
}

func loadSchema() Schema {
	var schema Schema
	//if file, err := os.Open("schema.json"); err != nil {
	//	fmt.Println(err)
	//	fmt.Println("failed to open schema")
	//	os.Exit(exitCodeUnexpected)
	//} else {
	//	defer file.Close()
	//	if schemajson, err := ioutil.ReadAll(file); err != nil {
	//		fmt.Println(err)
	//		fmt.Println("Failed to read schema")
	//		os.Exit(exitCodeUnexpected)
	if err := json.Unmarshal([]byte(schemaJson), &schema); err != nil {
		fmt.Println(err)
		fmt.Println("Invalid schema")
		os.Exit(exitCodeUnexpected)
	}
	//}
	return schema
}

func createCommandFromSchema(command SchemaCommand) *cobra.Command {
	var cmd *cobra.Command
	if command.Run.Cmd == "" {
		cmd = &cobra.Command{Use: command.Use, Short: command.Short, Long: command.Short}
	} else {
		cmd = &cobra.Command{
			Use: command.Use, Short: command.Short, Long: command.Short,
			Run: func(cmd *cobra.Command, args []string) {
				commandRun(cmd, command)
			},
		}
	}
	commandInit(cmd, command)
	return cmd
}
