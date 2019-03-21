package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

const cloudcliEnvPrefix = "CLOUDCLI"
const cloudcliConfigFilePrefix = ".cloudcli"

var enableAlpha = false

var rootCmd = &cobra.Command{
	Use:   "cloudcli",
	Short: "Cloudcli server management",
	Long: `Cloudcli server management - create, configure and manage servers`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		loadGlobalFlags()
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
	schema := loadSchema()
	for _, command := range schema.Commands {
		var cmd = createCommandFromSchema(command)
		for _, subcommand := range command.Commands {
			if ! subcommand.Alpha || enableAlpha {
				var subcmd = createCommandFromSchema(subcommand)
				cmd.AddCommand(subcmd)
			}
		}
		rootCmd.AddCommand(cmd)
	}
}

func initConfig() {
	if ! flagNoConfigValue {
		if flagConfigValue != "" {
			viper.SetConfigFile(flagConfigValue)
		} else if os.Getenv("CLOUDCLI_CONFIG") != "" {
			viper.SetConfigFile(os.Getenv("CLOUDCLI_CONFIG"))
		} else {
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			viper.AddConfigPath(home)
			viper.SetConfigName(cloudcliConfigFilePrefix)
		}
	}
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvPrefix(cloudcliEnvPrefix)
	if flagNoConfigValue {
		configFile = ""
	} else if err := viper.ReadInConfig(); err != nil {
		configFile = ""
	} else {
		configFile = viper.ConfigFileUsed()
	}
}
