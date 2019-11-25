package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

func returnGetCommandListResponse(outputFormat string, returnItems bool, resp_body []byte, command SchemaCommand, noExit bool, cmd *cobra.Command) []interface{} {
	var items []interface{}
	if outputFormat == "json" && ! returnItems {
		fmt.Println(string(resp_body))
		if ! noExit {
			os.Exit(0)
		}
	} else {
		if err := json.Unmarshal(resp_body, &items); err != nil {
			command_id, err := strconv.Atoi(string(resp_body))
			if err == nil {
				fmt.Println("Successfully queued command. Command ID:", command_id)
				waitForCommandIds(cmd, command, []string{strconv.Itoa(command_id)}, getCommandOutputFormat("", command, "human"), false);
				os.Exit(0)
			} else {
				fmt.Println(string(resp_body))
				os.Exit(exitCodeInvalidResponse)
			}
		}
		if ! returnItems {
			var outputItems []map[string]string;
			if command.Run.ParseStatisticsResponse {
				var metric string;
				var value string;
				var timestamp int64;
				for _, response := range items {
					for _, subResponse := range response.([]interface{}) {
						seriesResponse := subResponse.(map[string]interface{})
						metric = seriesResponse["series"].(string)
						for _, data := range seriesResponse["data"].([]interface{}) {
							dataItems := data.([]interface{})
							timestamp = int64(dataItems[0].(float64))
							value = parseItemString(dataItems[1])
							outputItem := make(map[string]string)
							outputItem["metric"] = metric
							outputItem["date"] = time.Unix(timestamp/1000, 0).Format("2006-01-02 15:04:05")
							outputItem["value"] = value
							outputItems = append(outputItems, outputItem)
						}
					}
				}
			} else {
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
							if ! noExit {
								os.Exit(0)
							}
						}
					}
				}
			}
			if outputFormat == "yaml" {
				if d, err := yaml.Marshal(&outputItems); err != nil {
					fmt.Println(string(resp_body))
					fmt.Println("Invalid response from server")
					os.Exit(exitCodeInvalidResponse)
				} else {
					fmt.Println(string(d))
					if ! noExit {
						os.Exit(0)
					}
				}
			} else {
				w := tabwriter.NewWriter(
					os.Stdout, 10, 0, 3, ' ',
					0,
				)
				var header []string
				for _, field := range command.Run.Fields {
					if ! field.Long && ! field.Hide {
						header = append(header, strings.ToUpper(field.Name))
					}
				}
				_, _ = fmt.Fprintf(w, "%s\n", strings.Join(header, "\t"))
				for _, outputItem := range outputItems {
					var row []string
					for _, field := range command.Run.Fields {
						if ! field.Hide {
							row = append(row, outputItem[field.Name])
						}
					}
					_, _ = fmt.Fprintf(w, "%s\n", strings.Join(row, "\t"))
				}
				_ = w.Flush()
				if ! noExit {
					os.Exit(0)
				}
			}
		}
	}
	return items
}

func commandRunGetList(cmd *cobra.Command, command SchemaCommand, returnItems bool, noWait bool, cmd_flags map[string]interface{}, outputFormat string, noExit bool) []interface{} {
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
			return commandRunGetListWaitFields(cmd, command, waitFields, cmd_flags, "human", noExit)
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
		if returnItems {
			return []interface{}{}
		} else {
			fmt.Println(resp.String())
			os.Exit(exitCodeInvalidStatus)
		}
	} else {
		items = returnGetCommandListResponse(outputFormat, returnItems, resp.Body(), command, noExit, cmd)
	}
	return items
}

