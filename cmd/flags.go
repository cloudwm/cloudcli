package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const flagApiServer = "api-server"
const flagApiServerKey = "apiServer"
var apiServer string
var flagApiServerValue string

const flagApiClientid = "api-clientid"
const flagApiClientidKey = "apiClientid"
var apiClientid string
var flagApiClientidValue string

const flagApiSecret = "api-secret"
const flagApiSecretKey = "apiSecret"
var apiSecret string
var flagApiSecretValue string

const flagConfig = "config"
var configFile string
var flagConfigValue string

const flagNoConfig = "no-config"
var flagNoConfigValue bool

const flagFormat = "format"
const flagFormatKey = "format"
var format string
var flagFormatValue string

const flagDebug = "debug"
const flagDebugKey = "debug"
var debug bool
var flagDebugValue bool

func loadGlobalFlags() {
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
}

func addGlobalFlags() {
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
}
