package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

func commandRunGetListWaitFields(cmd *cobra.Command, command SchemaCommand, waitFields []SchemaCommandField, cmd_flags map[string]interface{}, outputFormat string) []interface{} {
	var items []interface{}
	last_i := -1
	for {
		items = commandRunGetList(cmd, command, true, true, cmd_flags, outputFormat)
		var failed bool
		failed, last_i = printItemsCommandsProgress(items, waitFields, outputFormat, last_i)
		if failed {
			time.Sleep(5000000000)
		} else {
			break
		}
	}
	return commandRunGetList(cmd, command, false, true, cmd_flags, outputFormat)
}

func printItemsCommandsProgress(items []interface{}, waitFields []SchemaCommandField, outputFormat string, last_i int) (bool, int) {
	failed := false
	for _, item := range items {
		for _, field := range waitFields {
			fieldStatus := item.(map[string]interface{})[field.Name].(string)
			if fieldStatus != field.Wait {
				failed = true
			}
			last_i = printCommandProgress(items, field, outputFormat, item, last_i)
		}
	}
	return failed, last_i
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
