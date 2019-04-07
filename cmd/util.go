package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty"
	"gopkg.in/yaml.v2"
	"os"
)

func getJsonHttpResponse(path string) *resty.Response {
	var resp *resty.Response
	var err error
	if resp, err = resty.R().
		SetHeader("AuthClientId", apiClientid).
		SetHeader("AuthSecret", apiSecret).
		Get(fmt.Sprintf("%s%s", apiServer, path));
		err != nil {
		fmt.Println(err.Error())
		os.Exit(exitCodeUnexpected)
	} else if resp.StatusCode() != 200 {
		fmt.Println(resp.String())
		os.Exit(exitCodeInvalidStatus)
	}
	return resp
}

func jsonUnmarshalItemsList(respString string) map[string]interface{} {
	var items map[string]interface{}
	if err := json.Unmarshal([]byte(respString), &items); err != nil {
		fmt.Println(respString)
		fmt.Println("Invalid response from server")
		os.Exit(exitCodeInvalidResponse)
	}
	return items
}

func yamlDumpItemsList(respString string, items map[string]interface{}) {
	if d, err := yaml.Marshal(&items); err != nil {
		fmt.Println(respString)
		fmt.Println("Invalid response from server")
		os.Exit(exitCodeInvalidResponse)
	} else {
		fmt.Println(string(d))
	}
}


