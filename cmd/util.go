package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gopkg.in/resty.v1"
	"gopkg.in/yaml.v2"
	"os"
	"reflect"
	"strings"
)

func getJsonHttpResponse(path string) *resty.Response {
	var resp *resty.Response
	var err error
	if resp, err = resty.R().
		SetHeader("AuthClientId", apiClientid).
		SetHeader("AuthSecret", apiSecret).
		Get(fmt.Sprintf("%s%s", apiServer, path)); err != nil {
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

func parseItemString(item interface{}) string {
	var stringItem string
	switch typeditem := item.(type) {
	case float64:
		stringItem = fmt.Sprintf("%d", int(typeditem))
	case int:
		stringItem = fmt.Sprintf("%d", typeditem)
	case string:
		stringItem = fmt.Sprintf("%s", typeditem)
	case []interface{}:
		var stringSubItems []string
		for _, subitem := range typeditem {
			stringSubItems = append(stringSubItems, parseItemString(subitem))
		}
		stringItem = strings.Join(stringSubItems, ", ")
	case map[string]interface{}:
		var stringSubItems []string
		for k, v := range typeditem {
			stringSubItems = append(stringSubItems, fmt.Sprintf("%s=%s", k, parseItemString(v)))
		}
		stringItem = fmt.Sprintf("{%s}", strings.Join(stringSubItems, ", "))
	case bool:
		if typeditem {
			stringItem = "true"
		} else {
			stringItem = "false"
		}
	default:
		stringItem = fmt.Sprintf("unknown type: %s", reflect.TypeOf(item))
	}
	return stringItem
}

func getInput(title string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(title)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	text = strings.Replace(text, "\r", "", -1)
	text = strings.Replace(text, " ", "", -1)
	return text
}
