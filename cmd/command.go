package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"text/tabwriter"
)

func commandInit(cmd *cobra.Command, command SchemaCommand) {
	for _, flag := range command.Flags {
		cmd.Flags().StringP(flag.Name, "", "", flag.Usage)
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
	if resp, err := resty.R().
		SetHeader("AuthClientId", apiClientid).
		SetHeader("AuthSecret", apiSecret).
		Get(fmt.Sprintf("%s%s", apiServer, command.Run.Path));
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

func commandRunPost(cmd *cobra.Command, command SchemaCommand) {
	fmt.Println("post")
	os.Exit(exitCodeUnexpected)
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