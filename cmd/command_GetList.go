package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"
)

func commandRunGetList(cmd *cobra.Command, command SchemaCommand, returnItems bool, noWait bool, cmd_flags map[string]interface{}, outputFormat string) []interface{} {
	var qs []string
	var waitFields []SchemaCommandField
	for _, field := range command.Run.Fields {
		if field.Flag != "" {
			var value string;
			if field.Array {
				var arrayValue []string
				if cmd_flags != nil {
					arrayValue = cmd_flags[field.Flag].([]string)
				} else {
					arrayValue, _ = cmd.Flags().GetStringArray(field.Flag)
				}
				value = strings.Join(arrayValue, " ")
			} else if field.Bool || getSchemaCommandFlag(command, field.Flag).Bool {
				var boolValue bool
				if cmd_flags != nil {
					if cmd_flags[field.Flag] == nil {
						boolValue = false
					} else {
						boolValue = cmd_flags[field.Flag].(bool)
					}
				} else {
					boolValue, _ = cmd.Flags().GetBool(field.Flag)
				}
				if boolValue {
					value = "true"
				} else {
					value = ""
				}
			} else {
				if cmd_flags != nil {
					value = cmd_flags[field.Flag].(string)
				} else {
					value, _ = cmd.Flags().GetString(field.Flag)
				}
			}
			escapedValue := url.PathEscape(value)
			if (debug) {
				fmt.Printf("\nfield %s=%s / urlpart %s=%s", field.Flag, value, field.Name, escapedValue)
			}
			qs = append(qs, fmt.Sprintf("%s=%s", field.Name, escapedValue))
		} else if field.Wait != "" {
			waitFields = append(waitFields, field)
		}
	}
	outputFormat = getCommandOutputFormat(outputFormat, command, "human")
	if len(waitFields) > 0 && ! returnItems && ! noWait {
		var waitValue bool
		if cmd_flags != nil {
			waitValue = cmd_flags["wait"].(bool)
		} else {
			waitValue, _ = cmd.Flags().GetBool("wait");
		}
		if waitValue {
			return commandRunGetListWaitFields(cmd, command, waitFields, cmd_flags, outputFormat)
		}
	}
	var items []interface{}
	get_url := fmt.Sprintf("%s%s", apiServer, command.Run.Path)
	if len(qs) > 0 {
		get_url = fmt.Sprintf("%s?%s", get_url, strings.Join(qs, "&"))
	}
	if dryrun || debug {
		fmt.Printf("\nGET %s\n", get_url)
	}
	if dryrun {
		os.Exit(exitCodeDryrun)
	} else if resp, err := resty.R().
		SetHeader("AuthClientId", apiClientid).
		SetHeader("AuthSecret", apiSecret).
		Get(get_url);
		err != nil {
		fmt.Println(err.Error())
		os.Exit(exitCodeUnexpected)
	} else if resp.StatusCode() != 200 {
		fmt.Println(resp.String())
		os.Exit(exitCodeInvalidStatus)
	} else if outputFormat == "json" && ! returnItems {
		fmt.Println(resp.String())
		os.Exit(0)
	} else {
		if err := json.Unmarshal(resp.Body(), &items); err != nil {
			fmt.Println(resp.String())
			fmt.Println("Invalid response from server")
			os.Exit(exitCodeInvalidResponse)
		}
		if ! returnItems {
			var outputItems []map[string]string;
			for _, item := range items {
				inputItem := item.(map[string]interface{})
				outputItem := make(map[string]string)
				for _, field := range command.Run.Fields {
					if inputItem[field.Name] != nil {
						outputItem[field.Name] = parseItemString(inputItem[field.Name])
					}
				}
				outputItems = append(outputItems, outputItem)
			}
			if len(outputItems) == 1 && len(outputItems[0]) == 1 {
				for _, item := range outputItems {
					for _, v := range item {
						fmt.Println(v)
						os.Exit(0)
					}
				}
			}
			if outputFormat == "yaml" {
				if d, err := yaml.Marshal(&outputItems); err != nil {
					fmt.Println(resp.String())
					fmt.Println("Invalid response from server")
					os.Exit(exitCodeInvalidResponse)
				} else {
					fmt.Println(string(d))
					os.Exit(0)
				}
			} else {
				w := tabwriter.NewWriter(
					os.Stdout, 10, 0, 3, ' ',
					0,
				)
				var header []string
				for _, field := range command.Run.Fields {
					if ! field.Long {
						header = append(header, strings.ToUpper(field.Name))
					}
				}
				_, _ = fmt.Fprintf(w, "%s\n", strings.Join(header, "\t"))
				for _, outputItem := range outputItems {
					var row []string
					for _, field := range command.Run.Fields {
						row = append(row, outputItem[field.Name])
					}
					_, _ = fmt.Fprintf(w, "%s\n", strings.Join(row, "\t"))
				}
				_ = w.Flush()
				os.Exit(0)
			}
		}
	}
	return items
}

