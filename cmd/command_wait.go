package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"
)

func commandRunGetListWaitFields(cmd *cobra.Command, command SchemaCommand, waitFields []SchemaCommandField, cmd_flags map[string]interface{}, outputFormat string, noExit bool) []interface{} {
	var items []interface{}
	last_i := -1
	num_empty_responses := 0
	for {
		items = commandRunGetList(cmd, command, true, true, cmd_flags, outputFormat, false)
		if len(items) < 1 {
			if num_empty_responses > 10 {
				fmt.Println("Failed to get command status")
				os.Exit(exitCodeInvalidStatus)
			}
			time.Sleep(1000000000)
			num_empty_responses++
		} else {
			var failed bool
			var failedWithError bool
			failed, last_i, failedWithError = printItemsCommandsProgress(items, waitFields, outputFormat, last_i)
			if failedWithError {
				os.Exit(exitCodeInvalidStatus)
			} else if failed {
				time.Sleep(5000000000)
			} else {
				break
			}
		}
	}
	return commandRunGetList(cmd, command, false, true, cmd_flags, outputFormat, noExit)
}

func printItemsCommandsProgress(items []interface{}, waitFields []SchemaCommandField, outputFormat string, last_i int) (bool, int, bool) {
	failed := false
	failedWithError := false
	for _, item := range items {
		for _, field := range waitFields {
			fieldStatus := item.(map[string]interface{})[field.Name].(string)
			if fieldStatus != field.Wait {
				failed = true
			}
			if fieldStatus == field.WaitError {
				failedWithError = true
			}
			last_i = printCommandProgress(items, field, outputFormat, item, last_i)
		}
	}
	return failed, last_i, failedWithError
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


func waitForCommandIds(cmd *cobra.Command, command SchemaCommand, commandIds []string, outputFormat string, noExit bool) {
	if command.Wait {
		if b, _ := cmd.Flags().GetBool("wait"); b {
			s := ""
			if len(commandIds) > 1 {
				s = "s"
			}
			fmt.Printf("Waiting for command%s to complete\n", s)
			cmd_flags := make(map[string]interface{})
			cmd_flags["id"] = commandIds
			cmd_flags["wait"] = true
			_ = commandRunGetList(
				cli_cobra_subcommands["queue.detail"],
				cli_schema_subcommands["queue.detail"],
				false, false,
				cmd_flags, outputFormat, noExit,
			)
		}
	}
}
