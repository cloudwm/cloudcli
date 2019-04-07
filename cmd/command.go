package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"
)

func commandInit(cmd *cobra.Command, command SchemaCommand) {
	for _, flag := range command.Flags {
		if flag.Array {
			var defaultValue []string
			if flag.Default != "" {
				defaultValue = append(defaultValue, flag.Default)
			}
			cmd.Flags().StringArrayP(flag.Name, "", defaultValue, flag.Usage)
		} else if flag.Bool {
			cmd.Flags().BoolP(flag.Name, "", flag.Default != "", flag.Usage)
		} else {
			cmd.Flags().StringP(flag.Name, "", flag.Default, flag.Usage)
		}
		if flag.Required {
			cmd.MarkFlagRequired(flag.Name)
		}
	}
	if command.Run.Cmd == "getListOfLists" {
		commandInitGetListOfLists(cmd, command)
	}
}

func commandRun(cmd *cobra.Command, command SchemaCommand) {
	if command.Run.Cmd == "getList" {
		commandRunGetList(cmd, command)
	} else if command.Run.Cmd == "post" {
		commandRunPost(cmd, command)
	} else if command.Run.Cmd == "getListOfLists" {
		commandRunGetListOfLists(cmd, command)
	}
}

func commandRunGetList(cmd *cobra.Command, command SchemaCommand) {
	get_url := fmt.Sprintf("%s%s", apiServer, command.Run.Path)
	if dryrun {
		fmt.Printf("\nGET %s\n", get_url)
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
	} else if format == "json" {
		fmt.Println(resp.String())
		os.Exit(0)
	} else {
		var items []interface{}
		if err := json.Unmarshal(resp.Body(), &items); err != nil {
			fmt.Println(resp.String())
			fmt.Println("Invalid response from server")
			os.Exit(exitCodeInvalidResponse)
		}
		if format == "yaml" {
			if d, err := yaml.Marshal(&items); err != nil {
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
				header = append(header, strings.ToUpper(field.Name))
			}
			fmt.Fprintf(w, "%s\n", strings.Join(header, "\t"))
			for _, item := range items {
				var row []string
				for _, field := range command.Run.Fields {
					row = append(
						row,
						fmt.Sprintf("%s", item.(map[string]interface{})[field.Name]),
					)
				}
				fmt.Fprintf(w, "%s\n", strings.Join(row, "\t"))
			}
			w.Flush()
			os.Exit(0)
		}
	}
}

func commandExitErrorResponse(body []byte) {
	var errorResponse map[string]interface{}
	if err := json.Unmarshal(body, &errorResponse); err != nil {
		fmt.Println(string(body))
		fmt.Println("Failed to parse server error response")
		os.Exit(exitCodeInvalidResponse)
	} else {
		var message string;
		for k, v := range errorResponse {
			if k == "message" {
				message = v.(string);
			}
		}
		if format == "" && message != "" {
			fmt.Println(message)
		} else {
			var d []byte
			var err error
			if format == "yaml" {
				d, err = yaml.Marshal(&errorResponse)
			} else {
				d, err = json.Marshal(&errorResponse)
			}
			if err != nil {
				fmt.Println(string(body))
				fmt.Println("Invalid response from server")
				os.Exit(exitCodeInvalidResponse)
			} else {
				fmt.Println(string(d))
			}
		}
		os.Exit(exitCodeInvalidStatus)
	}
}

func commandRunPost(cmd *cobra.Command, command SchemaCommand) {
	var qs []string
	for _, field := range command.Run.Fields {
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
	if dryrun {
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
				commandExitErrorResponse(body)
			} else {
				var commandIds []string;
				if err := json.Unmarshal(body, &commandIds); err != nil {
					fmt.Println(string(body))
					fmt.Println("Failed to parse response")
					os.Exit(exitCodeInvalidResponse)
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
					os.Exit(0)
				} else {
					fmt.Println("Command IDs:")
					for _, commandId := range commandIds {
						fmt.Printf("%s\n", commandId)
					}
					os.Exit(0)
				}
			}
		}
	}
}

func commandInitGetListOfLists(cmd *cobra.Command, command SchemaCommand) {
	cmd.Flags().BoolP("cache", "", false,"save/load server options from file cloudcli-server-options.json")
	for _, list := range command.Run.Lists {
		cmd.Flags().BoolP(list.Name, "", false, fmt.Sprintf("only show %s resources", list.Name))
	}
}

func commandRunGetListOfLists(cmd *cobra.Command, command SchemaCommand) {
	var respString = ""
	var loadedFromCache = false
	cache, _ := cmd.Flags().GetBool("cache")
	if cache {
		if file, err := os.Open("cloudcli-server-options.json"); err == nil {
			defer file.Close()
			if respBytes, err := ioutil.ReadAll(file); err == nil {
				loadedFromCache = true
				respString = string(respBytes)
			}
		}
	}
	if ! loadedFromCache || respString == "" {
		respString = getJsonHttpResponse(command.Run.Path).String()
	}
	if cache && ! loadedFromCache {
		ioutil.WriteFile("cloudcli-server-options.json", []byte(respString), 0444)
	}
	onlyShow := ""
	for _, list := range command.Run.Lists {
		if b, err := cmd.Flags().GetBool(list.Name); err != nil {
			fmt.Println(err)
			os.Exit(exitCodeUnexpected)
		} else if b {
			if onlyShow != "" {
				fmt.Println("Invalid argumnets: only 1 show argument can be provided at a time")
			} else {
				onlyShow = list.Name
			}
		}
	}
	if format == "json" && onlyShow == "" {
		// json dump of all options
		fmt.Println(respString)
	} else {
		rootitems := jsonUnmarshalItemsList(respString)
		if format == "yaml" && onlyShow == "" {
			// yaml dump of all optoins
			yamlDumpItemsList(respString, rootitems)
		} else {
			var header []string
			if onlyShow == "" {
				fmt.Println("Not implemented yet, please choose an option (e.g. --datacenter / --cpu ..)")
				os.Exit(exitCodeUnexpected)
			} else {
				for _, list := range command.Run.Lists {
					if list.Name == onlyShow {
						header = getFieldsHeader(list.Fields, format)
						break
					}
				}
				if len(header) < 1 {
					fmt.Printf("Failed to get header for resource: %s\n", onlyShow)
					os.Exit(exitCodeUnexpected)
				}
			}
			w := tabwriter.NewWriter(os.Stdout, 10, 0, 3, ' ', 0)
			if format == "" {
				fmt.Fprintf(w, "%s\n", strings.Join(header, "\t"))
			}
			for rootlistkey, rootitem := range rootitems {
				for _, list := range command.Run.Lists {
					if list.Key == rootlistkey && (onlyShow == "" || onlyShow == list.Name) {
						if list.Type == "map" {
							printListKeyValueItemFields(rootitem.(map[string]interface{}), list.Fields, w, format)
						} else if list.Type == "list" {
							printListFields(rootitem.([]interface{}), list.Fields, w, format)
						} else if list.Type == "mapOfLists" {
							printMapOfListsFields(rootitem.(map[string]interface{}), list.Fields, w, format)
						} else {
							fmt.Printf("Invalid list type: %s\n%s\n", list.Type, list.Fields)
							os.Exit(exitCodeUnexpected)
						}
						break
					}
				}
			}
			if format == "" {
				w.Flush()
			}
		}
	}
}