package cmd

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"strings"
)

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

func parseItemString(item interface{}) string {
	var stringItem string
	switch typeditem := item.(type) {
	case float64:
		stringItem = fmt.Sprintf("%d", int(typeditem))
	default:
		stringItem = fmt.Sprintf("%s", typeditem)
	}
	return stringItem
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

func printListKeyValueItemFields(items map[string]interface{}, fields []SchemaCommandField, w io.Writer, format string) {
	if format == "" {
		for itemkey, itemvalue := range items {
			var row []string
			for _, value := range getKeyValueItemFields(fields, itemkey, itemvalue, format) {
				row = append(row, value)
			}
			fmt.Fprintf(w, "%s\n", strings.Join(row, "\t"))
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

func printListFields(items []interface{}, fields []SchemaCommandField, w io.Writer, format string) {
	if format == "" {
		for _, item := range items {
			var row []string
			for _, value := range getListItemFields(fields, item, format) {
				row = append(row, value)
			}
			fmt.Fprintf(w, "%s\n", strings.Join(row, "\t"))
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
				fmt.Fprintf(w, "%s\n", strings.Join(row, "\t"))
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
