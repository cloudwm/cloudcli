package cmd

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"text/tabwriter"

	"github.com/go-resty/resty"
	"github.com/spf13/cobra"
)

var serverListExitCodeUnexpected = 1
var serverListExitCodeInvalidStatus = 2
var serverListExitCodeInvalidResponse = 3

type ServerListItem struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Datacenter string `json:"datacenter"`
	Power string `json:"power"`
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List servers",
	Long: `List servers`,
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		if resp, err := resty.R().
			SetHeader("AuthClientId", apiClientid).
			SetHeader("AuthSecret", apiSecret).
			Get(fmt.Sprintf("%s/service/servers", apiServer));
		err != nil {
			fmt.Println(err.Error())
			os.Exit(serverListExitCodeUnexpected)
		} else if resp.StatusCode() != 200 {
			fmt.Println(resp.String())
			os.Exit(serverListExitCodeInvalidStatus)
		} else if format == "json" {
			fmt.Println(resp.String())
		} else {
			var servers []ServerListItem
			if err := json.Unmarshal(resp.Body(), &servers); err != nil {
				fmt.Println(resp.String())
				fmt.Println("Invalid response from server")
				os.Exit(serverListExitCodeInvalidResponse)
			}
			if format == "yaml" {
				if d, err := yaml.Marshal(&servers); err != nil {
					fmt.Println(resp.String())
					fmt.Println("Invalid response from server")
					os.Exit(serverListExitCodeInvalidResponse)
				} else {
					fmt.Println(string(d))
				}
			} else {
				w := tabwriter.NewWriter(
					os.Stdout, 10, 0, 3, ' ',
					0,
				)
				fmt.Fprintf(w, "ID\tNAME\tDATACENTER\tPOWER\n")
				for _, server := range servers {
					fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", server.Id, server.Name, server.Datacenter, server.Power)
				}
				w.Flush()
			}

		}
	},
}

func init() {
	serverCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
