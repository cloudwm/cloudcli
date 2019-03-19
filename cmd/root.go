package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var apiServer string
var flagApiServer = "api-server"
var flagApiServerValue string
var flagApiServerKey = "apiServer"

var apiClientid string
var flagApiClientid = "api-clientid"
var flagApiClientidValue string
var flagApiClientidKey = "apiClientid"

var apiSecret string
var flagApiSecret = "api-secret"
var flagApiSecretValue string
var flagApiSecretKey = "apiSecret"

var configFile string
var flagConfig = "config"
var flagConfigValue string

var noConfig bool
var flagNoConfig = "no-config"
var flagNoConfigValue bool

var format string
var flagFormat = "format"
var flagFormatValue string
var flagFormatKey = "format"

var debug bool
var flagDebug = "debug"
var flagDebugValue bool
var flagDebugKey = "debug"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cloudcli",
	Short: "Cloudcli server management",
	Long: `Cloudcli server management - create, configure and manage servers`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		debug = viper.GetBool(flagDebugKey)
		if debug {
			if configFile == "" {
				fmt.Fprintf(os.Stderr, "Not using config file\n")
			} else {
				fmt.Fprintf(os.Stderr, "Using config file: %s\n", configFile)
			}
		}
		numFailures := 0
		if apiServer = strings.TrimSpace(viper.GetString(flagApiServerKey)); apiServer == "" {
			numFailures += 1
			fmt.Printf("ERROR: --%s flag is required\n", flagApiServer)
		} else if debug {
			fmt.Fprintf(os.Stderr,"%s = %s\n", flagApiServerKey, apiServer)
		}
		if apiClientid = strings.TrimSpace(viper.GetString(flagApiClientidKey)); apiClientid == "" {
			numFailures += 1
			fmt.Printf("ERROR: --%s flag is required\n", flagApiClientid)
		} else if debug {
			fmt.Fprintf(os.Stderr,"%s = %s\n", flagApiClientidKey, apiClientid)
		}
		if apiSecret = strings.TrimSpace(viper.GetString(flagApiSecretKey)); apiSecret == "" {
			numFailures += 1
			fmt.Printf("ERROR: --%s flag is required\n", flagApiSecret)
		} else if debug {
			fmt.Fprintf(os.Stderr, "%s = %s\n", flagApiSecretKey, apiSecret)
		}
		format = strings.TrimSpace(viper.GetString(flagFormatKey))
		if format != "" && format != "json" && format != "yaml" {
			numFailures += 1
			fmt.Printf("ERROR: Unsupported --%s flag value: %s\n", flagFormatKey, format)
		} else if debug {
			fmt.Fprintf(os.Stderr, "%s = %s\n", flagFormatKey, format)
		}
		if numFailures > 0 {
			os.Exit(1)
		}
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&flagApiServerValue, flagApiServer, "", "API Server Hostname")
	viper.BindPFlag(flagApiServerKey, rootCmd.PersistentFlags().Lookup(flagApiServer))

	rootCmd.PersistentFlags().StringVar(&flagApiClientidValue, flagApiClientid, "", "API Client ID")
	viper.BindPFlag(flagApiClientidKey, rootCmd.PersistentFlags().Lookup(flagApiClientid))

	rootCmd.PersistentFlags().StringVar(&flagApiSecretValue, flagApiSecret, "", "API Secret")
	viper.BindPFlag(flagApiSecretKey, rootCmd.PersistentFlags().Lookup(flagApiSecret))

	rootCmd.PersistentFlags().StringVar(&flagConfigValue, flagConfig, "", "config file (default is $HOME/.cloudcli.yaml)")

	rootCmd.PersistentFlags().BoolVar(&flagNoConfigValue, flagNoConfig, false, "disable loading from config file")

	rootCmd.PersistentFlags().StringVar(&flagFormatValue, flagFormat, "", "output format, default format is a human readable summary")
	viper.BindPFlag(flagFormatKey, rootCmd.PersistentFlags().Lookup(flagFormat))

	rootCmd.PersistentFlags().BoolVar(&flagDebugValue, flagDebug, false, "enable debug output to stderr")
	viper.BindPFlag(flagDebugKey, rootCmd.PersistentFlags().Lookup(flagDebug))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if ! flagNoConfigValue {
		if flagConfigValue != "" {
			// Use config file from the flag.
			viper.SetConfigFile(flagConfigValue)
		} else {
			// Find home directory.
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// Search config in home directory with name ".cli" (without extension).
			viper.AddConfigPath(home)
			viper.SetConfigName(".cloudcli")
		}
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvPrefix("CLOUDCLI")

	if flagNoConfigValue {
		configFile = ""
	} else if err := viper.ReadInConfig(); err != nil {
		configFile = ""
	} else {
		configFile = viper.ConfigFileUsed()
	}
}
