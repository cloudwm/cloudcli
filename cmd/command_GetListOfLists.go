package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/tabwriter"
)

func commandInitGetListOfLists(cmd *cobra.Command, command SchemaCommand) {
	if command.CacheFile != "" {
		cmd.Flags().BoolP(
			"cache",
			"",
			false,
			fmt.Sprintf("save/load results from local file: %s", command.CacheFile),
		)
	}
	for _, list := range command.Run.Lists {
		cmd.Flags().BoolP(list.Name, "", false, fmt.Sprintf("only show %s resources", list.Name))
	}
}

func getListOfListsRespString(cacheFilePath string, enableCache bool, runPath string) string {
	var respString = ""
	var loadedFromCache = false
	cache := false
	if cacheFilePath != "" {
		cache = enableCache
		if cache {
			if file, err := os.Open(cacheFilePath); err == nil {
				defer file.Close()
				if respBytes, err := ioutil.ReadAll(file); err == nil {
					loadedFromCache = true
					respString = string(respBytes)
				}
			}
		}
	}
	if ! loadedFromCache || respString == "" {
		respString = getJsonHttpResponse(runPath).String()
	}
	if cache && ! loadedFromCache {
		_ = ioutil.WriteFile(cacheFilePath, []byte(respString), 0444)
	}
	return respString
}

func commandRunGetListOfLists(cmd *cobra.Command, command SchemaCommand) {
	enableCache, _ := cmd.Flags().GetBool("cache")
	respString := getListOfListsRespString(command.CacheFile, enableCache, command.Run.Path)
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
				_, _ = fmt.Fprintf(w, "%s\n", strings.Join(header, "\t"))
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
				_ = w.Flush()
			}
		}
	}
}


func printListKeyValueItemFields(items map[string]interface{}, fields []SchemaCommandField, w io.Writer, format string) {
	if format == "" {
		for itemkey, itemvalue := range items {
			var row []string
			for _, value := range getKeyValueItemFields(fields, itemkey, itemvalue, format) {
				row = append(row, value)
			}
			_, _ = fmt.Fprintf(w, "%s\n", strings.Join(row, "\t"))
		}
	} else {
		var dumpItems []map[string]string
		for itemkey, itemvalue := range items {
			dumpItem := make(map[string]string)
			for fieldnum, value := range getKeyValueItemFields(fields, itemkey, itemvalue, format) {
				field := fields[fieldnum]
				dumpItem[field.Name] = value
			}
			dumpItems = append(dumpItems, dumpItem)
		}
		if format == "yaml" {
			if d, err := yaml.Marshal(&dumpItems); err != nil {
				fmt.Println("Failed to create output yaml")
				os.Exit(exitCodeUnexpected)
			} else {
				fmt.Println(string(d))
			}
		} else {
			if d, err := json.Marshal(&dumpItems); err != nil {
				fmt.Println("Failed to create output json")
				os.Exit(exitCodeUnexpected)
			} else {
				fmt.Println(string(d))
			}
		}
	}
}


func printListFields(items []interface{}, fields []SchemaCommandField, w io.Writer, format string) {
	if format == "" {
		for _, item := range items {
			var row []string
			for _, value := range getListItemFields(fields, item, format) {
				row = append(row, value)
			}
			_, _ = fmt.Fprintf(w, "%s\n", strings.Join(row, "\t"))
		}
	} else {
		var dumpItems []map[string]string
		for _, item := range items {
			dumpItem := make(map[string]string)
			for fieldnum, value := range getListItemFields(fields, item, format) {
				field := fields[fieldnum]
				dumpItem[field.Name] = value
			}
			dumpItems = append(dumpItems, dumpItem)
		}
		if format == "yaml" {
			if d, err := yaml.Marshal(&dumpItems); err != nil {
				fmt.Println("Failed to create output yaml")
				os.Exit(exitCodeUnexpected)
			} else {
				fmt.Println(string(d))
			}
		} else {
			if d, err := json.Marshal(&dumpItems); err != nil {
				fmt.Println("Failed to create output json")
				os.Exit(exitCodeUnexpected)
			} else {
				fmt.Println(string(d))
			}
		}
	}
}


func printMapOfListsFields(itemsmap map[string]interface{}, fields []SchemaCommandField, w io.Writer, format string) {
	if format == "" {
		for itemsmapkey, itemsmapvalue := range itemsmap {
			for _, item := range itemsmapvalue.([]interface{}) {
				var row []string
				for _, value := range getKeyValueItemFields(fields, itemsmapkey, item, format) {
					row = append(row, value)
				}
				_, _ = fmt.Fprintf(w, "%s\n", strings.Join(row, "\t"))
			}
		}
	} else {
		var dumpItems []map[string]string
		for itemsmapkey, itemsmapvalue := range itemsmap {
			for _, item := range itemsmapvalue.([]interface{}) {
				dumpItem := make(map[string]string)
				for fieldnum, value := range getKeyValueItemFields(fields, itemsmapkey, item, format) {
					field := fields[fieldnum]
					dumpItem[field.Name] = value
				}
				dumpItems = append(dumpItems, dumpItem)
			}
		}
		if format == "yaml" {
			if d, err := yaml.Marshal(&dumpItems); err != nil {
				fmt.Println("Failed to create output yaml")
				os.Exit(exitCodeUnexpected)
			} else {
				fmt.Println(string(d))
			}
		} else {
			if d, err := json.Marshal(&dumpItems); err != nil {
				fmt.Println("Failed to create output json")
				os.Exit(exitCodeUnexpected)
			} else {
				fmt.Println(string(d))
			}
		}
	}
}


func getListItemFields(fields []SchemaCommandField, item interface{}, format string) []string {
	stringItem := parseItemString(item)
	var values []string
	for _, field := range fields {
		value := ""
		if field.From == "value" {
			value = stringItem
		}
		value = runFieldParsers(field, value, format)
		values = append(values, value)
	}
	return values
}


func getKeyValueItemFields(fields []SchemaCommandField, itemkey interface{}, itemvalue interface{}, format string) []string {
	var values []string
	for _, field := range fields {
		if field.Long && format == "" {
			continue
		}
		value := ""
		if field.From == "key" {
			value = parseItemString(itemkey)
		} else if field.From == "value" {
			value = parseItemString(itemvalue)
		} else {
			for k, v := range itemvalue.(map[string]interface{}) {
				if k == field.From {
					value = parseItemString(v)
					break
				}
			}
		}
		value = runFieldParsers(field, value, format)
		values = append(values, value)
	}
	return values
}


func runFieldParsers(field SchemaCommandField, value string, format string) string {
	for _, parser := range field.Parsers {
		if parser.OnlyForHumans && format != "" {
			continue
		}
		if parser.Parser == "split_value_remove_first" {
			splitString := parser.SplitString
			value = strings.TrimSpace(strings.Join(
				strings.Split(value, splitString)[1:], splitString,
			))
		} else if parser.Parser == "network_ips" {
			value = strings.TrimSpace(value)
			if strings.Contains(value, " ") {
				value = fmt.Sprintf("%d", len(strings.Split(value, " ")))
			}
		}
	}
	return value
}


func getFieldsHeader(fields []SchemaCommandField, format string) []string {
	var header []string
	for _, field := range fields {
		if field.Long && format == "" {
			continue
		}
		header = append(header, strings.ToUpper(field.Name))
	}
	return header
}
