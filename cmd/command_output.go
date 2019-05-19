package cmd

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

func getCommandOutputFormat(outputFormat string, command SchemaCommand, defaultOutputFormat string) string {
	if command.DefaultFormat != "" {
		defaultOutputFormat = command.DefaultFormat
	}
	return getOutputFormat(outputFormat, defaultOutputFormat)
}


func commandExitErrorResponse(body []byte, command SchemaCommand) {
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
			fmt.Printf("%s command failed: %s\n", command.Use, message)
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
			} else if format == "" {
				fmt.Printf("%s command failed: %s\n", command.Use, string(d))
			} else {
				fmt.Println(string(d))
			}
		}
		os.Exit(exitCodeInvalidStatus)
	}
}
