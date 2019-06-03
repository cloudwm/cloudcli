package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type ServerIdsDryrunResponse struct {
	Dryrun bool `json:"dryrun"`
	ServerNames []string `json:"server-names"`
}

func commandRunPost(cmd *cobra.Command, command SchemaCommand) {
	var qs []string
	hasDryrunFlag := false
	for _, field := range command.Run.Fields {
		if field.Flag == "dryrun" {
			hasDryrunFlag = true
			if dryrun {
				qs = append(qs, "dryrun=true")
			}
			continue
		}
		var value string;
		if field.Array {
			arrayValue, _ := cmd.Flags().GetStringArray(field.Flag)
			value = strings.Join(arrayValue, " ")
		} else if field.Bool {
			if boolValue, _ := cmd.Flags().GetBool(field.Flag); boolValue {
				value = "true"
			} else {
				value = ""
			}
		} else {
			value, _ = cmd.Flags().GetString(field.Flag)
		}
		escapedValue := url.PathEscape(value)
		if (debug) {
			fmt.Printf("\nfield %s=%s / urlpart %s=%s", field.Flag, value, field.Name, escapedValue)
		}
		qs = append(qs, fmt.Sprintf("%s=%s", field.Name, escapedValue))
	}
	payload := strings.Join(qs, "&")
	post_url := fmt.Sprintf("%s%s", apiServer, command.Run.Path)
	if dryrun && ! hasDryrunFlag {
		fmt.Printf("\nPOST %s\n", post_url)
		fmt.Printf("%s\n\n", payload)
		os.Exit(exitCodeDryrun)
	} else {
		if req, err := http.NewRequest("POST", post_url, strings.NewReader(payload)); err != nil {
			fmt.Println("Failed to create POST request")
			os.Exit(exitCodeUnexpected)
		} else {
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Add("AuthClientId", apiClientid)
			req.Header.Add("AuthSecret", apiSecret)
			if r, err := http.DefaultClient.Do(req); err != nil {
				fmt.Println("Failed to send POST request")
				os.Exit(exitCodeUnexpected)
			} else if body, err := ioutil.ReadAll(r.Body); err != nil {
				fmt.Println("Failed to read POST response body")
				os.Exit(exitCodeInvalidResponse)
			} else if r.StatusCode != 200 {
				commandExitErrorResponse(body, command)
			} else if command.Run.ServerMethod == "GET" {
				returnGetCommandListResponse(
					getCommandOutputFormat("", command, "human"),
					false, body, command,
				)
			} else {
				var commandIds []string;
				if err := json.Unmarshal(body, &commandIds); err != nil {
					if dryrun && hasDryrunFlag {
						var dryrunRes ServerIdsDryrunResponse
						if err := json.Unmarshal(body, &dryrunRes); err != nil {
							fmt.Println(string(body))
							fmt.Println("Failed to parse dryrun response")
							os.Exit(exitCodeInvalidResponse)
						} else if format == "yaml" || format == "json" {
							var d []byte
							var err error
							if format == "yaml" {
								d, err = yaml.Marshal(&dryrunRes)
							} else {
								d, err = json.Marshal(&dryrunRes)
							}
							if err != nil {
								fmt.Println(string(body))
								fmt.Println("Invalid response from server")
								os.Exit(exitCodeInvalidResponse)
							} else {
								fmt.Println(string(d))
								os.Exit(0)
							}
						} else {
							fmt.Printf("server names to delete: %s\n", strings.Join(dryrunRes.ServerNames, ", "))
							os.Exit(0)
						}
					} else {
						fmt.Println(string(body))
						fmt.Println("Failed to parse response")
						os.Exit(exitCodeInvalidResponse)
					}
				}
				if len(commandIds) == 0 {
					fmt.Println("Unexpected command failure")
					os.Exit(exitCodeUnexpected);
				}
				if format == "json" || format == "yaml" {
					parsedResponse := make(map[string][]string);
					parsedResponse["command_ids"] = commandIds;
					var d []byte
					var err error
					if format == "yaml" {
						d, err = yaml.Marshal(&parsedResponse)
					} else {
						d, err = json.Marshal(&parsedResponse)
					}
					if err != nil {
						fmt.Println(string(body))
						fmt.Println("Invalid response from server")
						os.Exit(exitCodeInvalidResponse)
					} else {
						fmt.Println(string(d))
						os.Exit(0)
					}
				} else if len(commandIds) == 1 {
					fmt.Printf("Command ID: %s\n", commandIds[0])
					waitForCommandIds(cmd, command, commandIds, getCommandOutputFormat("", command, "human"));
					os.Exit(0)
				} else {
					fmt.Println("Command IDs:")
					for _, commandId := range commandIds {
						fmt.Printf("%s\n", commandId)
					}
					waitForCommandIds(cmd, command, commandIds, getCommandOutputFormat("", command, "human"))
					os.Exit(0)
				}
			}
		}
	}
}
