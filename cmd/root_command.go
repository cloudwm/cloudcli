package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

const cloudcliEnvPrefix = "CLOUDCLI"
const cloudcliConfigFilePrefix = ".cloudcli"

var enableAlpha = false
var hasSchema = false
var cli_schema Schema

var cli_cobra_subcommands map[string]*cobra.Command
var cli_schema_subcommands map[string]SchemaCommand

func getRootLongDescription() string {
	return `Cloudcli server management - create, configure and manage servers

## Configuring cloudcli

You can set credentials and arguments using one of the following options:

* A yaml configuration file at HOME/.cloudcli.yaml (or specified using the CLOUDCLI_CONFIG env var or --config flag)
* Environment variables - uppercase strings, split with underscore and prefixed with CLOUDCLI_
* CLI flags: --api-server "" --api-clientid "" --api-secret ""

See [example.cloudcli.yaml](https://github.com/cloudwm/cloudcli/blob/master/example-cloudcli.yaml) and [example-cloudcli.env](https://github.com/cloudwm/cloudcli/blob/master/example-cloudcli.env) for more details on using the yaml config file or environment variables.

**Important** Please keep your server and API credentials secure, 
it's recommended to use a configuration file with appropriate permissions and location.
`
}

var rootCmd = &cobra.Command{
	Use:   "cloudcli",
	Short: "Cloudcli server management",
	Long: getRootLongDescription(),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		loadGlobalFlags()
	},
}

var initCommand = getInitCommand()

var versionCommand = &cobra.Command{
	Use: "version",
	Short: "Get the cloudcli version",
	Run: func(cmd *cobra.Command, args []string) {
		var schemaVersion []string
		for _, versionPart := range cli_schema.SchemaVersion {
			schemaVersion = append(schemaVersion, fmt.Sprintf("%d", versionPart))
		}
		outputFormat := getOutputFormat("", "human")
		if outputFormat == "human" {
			versionString := strings.Join(schemaVersion, ".")
			fmt.Printf("cloudcli v%s\n", versionString)
			os.Exit(0)
		} else if outputFormat == "json" {
			versionString := strings.Join(schemaVersion, ", ")
			fmt.Printf("{\"cloudcli-version\": [%s]}\n", versionString)
			os.Exit(0)
		} else if outputFormat == "yaml" {
			versionString := strings.Join(schemaVersion, ", ")
			fmt.Printf("cloudcli-version: [%s]\n", versionString)
			os.Exit(0)
		} else {
			fmt.Println("Invalid output format")
			os.Exit(exitCodeUnexpected)
		}
	},
}


func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	addGlobalFlags()
	if os.Getenv("CLOUDCLI_ENABLE_ALPHA") != "" {
		enableAlpha = true
	}
	rootCmd.SetUsageTemplate(`{{if .HasAvailableInheritedFlags}}Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}

{{end}}Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)
	rootCmd.AddCommand(initCommand)
	rootCmd.AddCommand(versionCommand)
	cli_cobra_subcommands = make(map[string]*cobra.Command)
	cli_schema_subcommands = make(map[string]SchemaCommand)
	hasSchema, cli_schema = loadSchema()
	if ! hasSchema {
		_ = rootCmd.PersistentFlags().Parse(os.Args[1:])
		initConfig()
		loadGlobalFlags()
		cli_schema = downloadSchema(schemaFile, fmt.Sprintf("%s%s", apiServer, "/schema"))
		hasSchema = true
	}
	initSubCommands()
}

func initSubCommands() {
	for _, command := range cli_schema.Commands {
		var cmd= createCommandFromSchema(command)
		for _, subcommand := range command.Commands {
			if ! subcommand.Alpha || enableAlpha {
				var subcmd= createCommandFromSchema(subcommand)
				cmd.AddCommand(subcmd)
				cli_cobra_subcommands[fmt.Sprintf("%s.%s", command.Use, subcommand.Use)] = cmd
				cli_schema_subcommands[fmt.Sprintf("%s.%s", command.Use, subcommand.Use)] = subcommand
			}
		}
		rootCmd.AddCommand(cmd)
	}
}
