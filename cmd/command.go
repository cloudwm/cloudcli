package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func commandInit(cmd *cobra.Command, command SchemaCommand) {
	for _, flag := range command.Flags {
		if flag.Name == "dryrun" {
			// this flag is already added as a global flag
			continue
		}
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
			_ = cmd.MarkFlagRequired(flag.Name)
		}
	}
	if command.Run.Cmd == "getListOfLists" {
		commandInitGetListOfLists(cmd, command)
	}
}

func commandRun(cmd *cobra.Command, command SchemaCommand) {
	for _, flag := range command.Flags {
		for _, processing := range flag.Processing {
			if processing.Method == "validateAllowedOutputFormats" {
				if ! flag.Bool {
					fmt.Println("Unexpected error (validate output formats requires bool flag)")
					os.Exit(exitCodeInvalidFlags)
				}
				if b, _ := cmd.Flags().GetBool(flag.Name); b {
					outputFormat := getCommandOutputFormat("", command, "human")
					ok := false
					allowedOutputFormats := processing.Args.([]interface{})
					for _, allowedFormat := range allowedOutputFormats {
						if allowedFormat.(string) == outputFormat {
							ok = true
						}
					}
					if ! ok {
						fmt.Printf("%s flag can't be used with %s format\n", flag.Name, outputFormat)
						os.Exit(exitCodeInvalidFlags)
					}
				}
			}
		}
	}
	numErrors := 0
	if apiClientid == "" {
		fmt.Printf("ERROR: --%s flag is required\n", flagApiClientid)
		numErrors += 1
	}
	if apiSecret == "" {
		fmt.Printf("ERROR: --%s flag is required\n", flagApiSecret)
		numErrors += 1
	}
	if numErrors > 0 {
		os.Exit(exitCodeInvalidFlags)
	}
	if command.Run.Cmd == "getList" {
		commandRunGetList(cmd, command, false, false, nil, "")
	} else if command.Run.Cmd == "post" {
		commandRunPost(cmd, command)
	} else if command.Run.Cmd == "getListOfLists" {
		commandRunGetListOfLists(cmd, command)
	}
}


