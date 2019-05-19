package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

func commandRunGetListWaitFields(cmd *cobra.Command, command SchemaCommand, waitFields []SchemaCommandField, cmd_flags map[string]interface{}, outputFormat string) []interface{} {
	var items []interface{}
	for {
		items = commandRunGetList(cmd, command, true, true, cmd_flags, outputFormat)
		if printItemsCommandsProgress(items, waitFields, outputFormat) {
			time.Sleep(5000000000)
		} else {
			break
		}
	}
	return commandRunGetList(cmd, command, false, true, cmd_flags, outputFormat)
}

func printItemsCommandsProgress(items []interface{}, waitFields []SchemaCommandField, outputFormat string) bool {
	failed := false
	last_i := -1
	for _, item := range items {
		for _, field := range waitFields {
			fieldStatus := item.(map[string]interface{})[field.Name].(string)
			if fieldStatus != field.Wait {
				failed = true
			}
			last_i = printCommandProgress(items, field, outputFormat, item, last_i)
		}
	}
	return failed
}

func printCommandProgress(items []interface{}, field SchemaCommandField, outputFormat string, item interface{}, last_i int) int {
	if len(items) == 1 && field.WaitPrintField != "" && outputFormat == "human" {
		for i, line := range strings.Split(item.(map[string]interface{})[field.WaitPrintField].(string), "\n") {
			if i > last_i {
				last_i = i
				fmt.Printf("%s\n", line)
			}
		}
	}
	// else {
	// TODO: figure out what to output when progressing on multiple command ids
	// at the moment it affects terminate / power operations which complete quickly
	//fmt.Printf("%s=%s\n", field.Name, fieldStatus)
	//}
	return last_i
}


func waitForCommandIds(cmd *cobra.Command, command SchemaCommand, commandIds []string, outputFormat string) {
	if command.Wait {
		if b, _ := cmd.Flags().GetBool("wait"); b {
			fmt.Println("Waiting for commands to complete")
			cmd_flags := make(map[string]interface{})
			cmd_flags["id"] = commandIds
			cmd_flags["wait"] = true
			_ = commandRunGetList(
				cli_cobra_subcommands["queue.detail"],
				cli_schema_subcommands["queue.detail"],
				false, false,
				cmd_flags, outputFormat,
			)
		}
	}
}
