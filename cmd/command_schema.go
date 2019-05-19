package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"time"
)

var schemaFile string

type SchemaCommandFieldParser struct {
	Parser string `json:"parser"`
	SplitString string `json:"split_string"`
	OnlyForHumans bool `json:"only_for_humans"`
}

type SchemaCommandField struct {
	Name string `json:"name"`
	Flag string `json:"flag"`
	From string `json:"from"`
	Parsers []SchemaCommandFieldParser `json:"parsers"`
	Long bool `json:"long"`
	Array bool `json:"array"`
	Bool bool `json:"bool"`
	Number bool `json:"number"`
	Wait string `json:"wait"`
	WaitPrintField string `json:"waitPrintField"`
}

type SchemaCommandFlagProcessing struct {
	Method string `json:"method"`
	Args interface{} `json:"args"`
}

type SchemaCommandFlag struct {
	Name string `json:"name"`
	Usage string `json:"usage"`
	Required bool `json:"required"`
	Array bool `json:"array"`
	Default string `json:"default"`
	Bool bool `json:"bool"`
	Processing []SchemaCommandFlagProcessing `json:"processing"`
}

type SchemaCommandList struct {
	Name string `json:"name"`
	Key string `json:"key"`
	Type string `json:"type"`
	Fields []SchemaCommandField `json:"fields"`

}

type SchemaCommandRun struct {
	Cmd string `json:"cmd"`
	Path string `json:"path"`
	Fields []SchemaCommandField `json:"fields"`
	Lists []SchemaCommandList `json:"lists"`
}

type SchemaCommand struct {
	Alpha bool `json:"alpha"`
	Use string `json:"use"`
	Short string `json:"short"`
	Run SchemaCommandRun `json:"run"`
	Flags []SchemaCommandFlag `json:"flags"`
	Commands []SchemaCommand `json:"commands"`
	Wait bool `json:"wait"`
	DefaultFormat string `json:"default-format"`
	CacheFile string `json:"cache-file"`
}

type Schema struct {
	SchemaVersion [3]int `json:"schema_version"`
	Commands []SchemaCommand `json:"commands"`
	SchemaGeneratedAt time.Time `json:"schema_generated_at"`
}

func downloadSchema(schemaFile string, schemaUrl string) Schema {
	var schema_ Schema
	if dryrun || debug {
		if debug {
			_, _ = fmt.Fprintf(os.Stderr, "\nAuthClientId: %s", apiClientid)
			_, _ = fmt.Fprintf(os.Stderr, "\nAuthSecret: %s", apiSecret)
		}
		_, _ = fmt.Fprintf(os.Stderr, "\nGET %s\n", schemaUrl)
	}
	if dryrun {
		os.Exit(exitCodeDryrun)
	} else if resp, err := resty.R().
		SetHeader("AuthClientId", apiClientid).
		SetHeader("AuthSecret", apiSecret).
		Get(schemaUrl); err != nil {
		fmt.Println(err.Error())
		os.Exit(exitCodeUnexpected)
	} else if resp.StatusCode() != 200 {
		fmt.Println(resp.String())
		os.Exit(exitCodeInvalidStatus)
	} else if format == "json" {
		fmt.Println(resp.String())
		os.Exit(0)
	} else if err := ioutil.WriteFile(schemaFile, []byte(resp.String()), 0644); err != nil {
		fmt.Printf("Failed  to write schemaFile (%s)", schemaFile)
		os.Exit(exitCodeUnexpected)
	} else if err := json.Unmarshal([]byte(resp.String()), &schema_); err != nil {
		fmt.Println(err)
		fmt.Println("Invalid schema")
		os.Exit(exitCodeUnexpected)
	}
	return schema_
}

func loadSchema() (bool, Schema) {
	var schema Schema
	has_schema := false
	if schemaFile = os.Getenv("CLOUDCLI_SCHEMA_FILE"); schemaFile == "" {
		if home, err := homedir.Dir(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			schemaFile = fmt.Sprintf("%s/%s", home, ".cloudcli.schema.json")
		}
	}
	if file, err := os.Open(schemaFile); err == nil {
		defer file.Close()
		if schemaJsonString, err := ioutil.ReadAll(file); err != nil {
			fmt.Println(err)
			fmt.Println("Failed to read schema")
			_ = os.Remove(schemaFile)
		} else if err := json.Unmarshal([]byte(schemaJsonString), &schema); err != nil {
			fmt.Println(err)
			fmt.Println("Invalid schema")
			_ = os.Remove(schemaFile)
		} else if time.Since(schema.SchemaGeneratedAt).Minutes() < 60 {
			has_schema = true
		}
	}
	return has_schema, schema
}

func createCommandFromSchema(command SchemaCommand) *cobra.Command {
	var cmd *cobra.Command
	if command.Run.Cmd == "" {
		cmd = &cobra.Command{Use: command.Use, Short: command.Short, Long: command.Short}
	} else {
		cmd = &cobra.Command{
			Use: command.Use, Short: command.Short, Long: command.Short,
			Run: func(cmd *cobra.Command, args []string) {
				commandRun(cmd, command)
			},
		}
	}
	commandInit(cmd, command)
	return cmd
}

func getSchemaCommandFlag(command SchemaCommand, flagName string) SchemaCommandFlag {
	for _, flag := range command.Flags {
		if flagName == flag.Name {
			return flag
		}
	}
	return SchemaCommandFlag{}
}
