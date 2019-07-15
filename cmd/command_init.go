package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"os"
)

func getInitCommand() *cobra.Command {
	cmd := cobra.Command{
		Use: "init",
		Short: "Initialize cloudcli",
		Long: "Authenticates to a cloudcli server and updates CLI to latest version",
		Run: func(cmd *cobra.Command, args []string) {
			schema := downloadSchema(schemaFile, fmt.Sprintf("%s%s", apiServer, "/schema"))
			if apiClientid == "" && apiSecret == "" {
				fmt.Printf("Performing interactive initialization of cloudcli\n")
				apiClientid = getInput("Enter the API Client ID: ")
				apiSecret = getInput("Enter the API Secret: ")
				if apiClientid == "" || apiSecret == "" {
					fmt.Println("Missing API client ID / Secret")
					os.Exit(exitCodeUnexpected)
				}
				_ = writeNewConfigFile(fmt.Sprintf("apiServer: \"%s\"\napiClientid: \"%s\"\napiSecret: \"%s\"\n", apiServer, apiClientid, apiSecret))
			}
			fmt.Printf(
				"cloudcli v%d.%d.%d Initialized successfully.\n\nYou can now run cloudcli commands, see:\ncloudcli --help\n\n",
				schema.SchemaVersion[0], schema.SchemaVersion[1], schema.SchemaVersion[2],
			)
			docsDir, _ := cmd.Flags().GetString("docs-dir")
			docsFormat, _ := cmd.Flags().GetString("docs-format")
			if docsFormat == "markdown" {
				file_info, err := os.Stat(docsDir)
				if err != nil {
					fmt.Print(err, "\n")
					os.Exit(exitCodeUnexpected)
				}
				if ! file_info.IsDir() {
					fmt.Printf("Please create the docs directory at %s\n", docsDir)
					os.Exit(exitCodeUnexpected)
				}
				err = doc.GenMarkdownTree(rootCmd, docsDir)
				if err != nil {
					fmt.Print(err, "\n")
					os.Exit(exitCodeUnexpected)
				}
				fmt.Printf("Successfully created %s docs in directory %s\n", docsFormat, docsDir)
				os.Exit(0)
			} else if docsFormat != "" {
				fmt.Printf("Invalid docs format: %s\n", docsFormat)
				os.Exit(exitCodeInvalidFlags)
			}
		},
	}
	cmd.Flags().StringP("docs-dir", "", "docs", "Save generated docs to this directory")
	cmd.Flags().StringP("docs-format", "", "", "Set to one of the supported formats to generate docs. supported formats: markdown")
	return &cmd
}
