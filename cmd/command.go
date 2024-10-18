package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func commandInit(cmd *cobra.Command, command SchemaCommand) {
	if command.DontSortFlags {
		cmd.Flags().SortFlags = false
	}
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
		// to support interactive mode - we conditionally set flags as required only if not in interactive mode
		// check of required flag is done server-side
		// TODO: determine a way to mark as required but also support server create --interactive
		//if flag.Required {
		//	_ = cmd.MarkFlagRequired(flag.Name)
		//}
	}
	if command.Run.Cmd == "getListOfLists" {
		commandInitGetListOfLists(cmd, command)
	}
	cliUsage := command.CliUsage
	if cliUsage != "" {
		cliUsage = fmt.Sprintf("\n\n%s", cliUsage)
	}
	cmd.SetUsageTemplate(fmt.Sprintf(`Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}%s{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Use "cloudcli --help" for a list of available global flags and general usage instructions.{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`, cliUsage))
}

func commandPreRun(cmd *cobra.Command, command SchemaCommand) {
	for _, hook := range command.CliPreRunHooks {
		if hook.Type == "requireOneOf" {
			numFlags := 0
			for _, name := range hook.OneOf {
				val, _ := cmd.Flags().GetString(name)
				if val != "" {
					numFlags += 1
				}
			}
			if numFlags != 1 {
				fmt.Printf(
					"syntax error, invalid arguments, use cloudcli %s %s --help for more information\n",
					cmd.Parent().Use, cmd.Use,
				)
				os.Exit(exitCodeInvalidFlags)
			}
		}
	}
}

func commandRun(cmd *cobra.Command, command SchemaCommand) {
	for _, flag := range command.Flags {
		for _, processing := range flag.Processing {
			if processing.Method == "validateAllowedOutputFormats" {
				if !flag.Bool {
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
					if !ok {
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
		fmt.Println("Provide the missing flags or run `cloudcli init` to initialize interactively")
		os.Exit(exitCodeInvalidFlags)
	}
	if command.Run.Cmd == "getList" {
		commandRunGetList(cmd, command, false, false, nil, "", false)
	} else if command.Run.Cmd == "post" {
		commandRunPost(cmd, command)
	} else if command.Run.Cmd == "getListOfLists" {
		commandRunGetListOfLists(cmd, command)
	}
}
