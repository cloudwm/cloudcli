package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"io"
	"os"
	"path"
)

func initConfig() {
	configFilePath := getConfigFilePath()
	if configFilePath != "" {
		viper.SetConfigFile(configFilePath)
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

func getConfigFilePath() string {
	if flagNoConfigValue {
		return ""
	} else {
		if flagConfigValue != "" {
			return flagConfigValue
		} else if os.Getenv("CLOUDCLI_CONFIG") != "" {
			return os.Getenv("CLOUDCLI_CONFIG")
		} else {
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return path.Join(home, cloudcliConfigFilePrefix + ".yaml")
		}
	}
}

func writeNewConfigFile(content string) error {
	filename := getConfigFilePath()
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.WriteString(file, content)
	if err != nil {
		return err
	}
	return file.Sync()
}
