package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type ServersSshInfo struct {
	ExternalIp string `json:"externalIp"`
}

func getServerIP(serverName string) string {
	var items []interface{}
	post_url := fmt.Sprintf("%s%s", apiServer, "/service/server/info")
	payload := fmt.Sprintf("name=%s", serverName)
	if req, err := http.NewRequest("POST", post_url, strings.NewReader(payload)); err != nil {
		fmt.Println("Failed to create POST request")
		os.Exit(exitCodeUnexpected)
	} else {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("AuthClientId", apiClientid)
		req.Header.Add("AuthSecret", apiSecret)
		if r, err := http.DefaultClient.Do(req); err != nil {
			fmt.Println("Failed to send POST request")
			os.Exit(exitCodeUnexpected)
		} else if body, err := ioutil.ReadAll(r.Body); err != nil {
			fmt.Println("Failed to read POST response body")
			os.Exit(exitCodeInvalidResponse)
		} else if r.StatusCode != 200 {
			fmt.Printf("Failed to get server IP\n")
			os.Exit(exitCodeUnexpected)
		} else if err := json.Unmarshal(body, &items); err != nil {
			fmt.Println(string(body))
			fmt.Println("Invalid response from server")
			os.Exit(exitCodeInvalidResponse)
		} else if len(items) != 1 {
			fmt.Printf("Did not find matching server\n")
			os.Exit(exitCodeUnexpected)
		} else {
			return items[0].(map[string]interface{})["networks"].([]interface{})[0].(map[string]interface{})["ips"].([]interface{})[0].(string)
		}
	}
	return ""
}

func setServerSshKey(serverPassword string, serverIp string, publicKey string) bool {
	fmt.Printf("Setting SSH key from %s to server IP %s\n", publicKey, serverIp)
	server := serverIp + ":22"
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password(serverPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", server, config)
	if err != nil {
		fmt.Printf("Failed to dial: %s\n", err.Error())
		return false
	}
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil {
		fmt.Printf("Failed to create session: %s\n", err.Error())
		return false
	}
	defer session.Close()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	publicKeyBytes, err := ioutil.ReadFile(publicKey)
	if err != nil {
		fmt.Printf("Failed to read public key file\n")
		os.Exit(exitCodeUnexpected)
	}
	err = session.Run(fmt.Sprintf("echo '%s' >> ~/.ssh/authorized_keys", string(publicKeyBytes)))
	if err != nil {
		fmt.Printf("Failed to add public key to authorized keys\n")
		os.Exit(exitCodeUnexpected)
	}
	fmt.Printf("Successfuly added public key to authorized keys on the server\n")
	return true
}

func commandRunSsh(cmd *cobra.Command, command SchemaCommand, serversInfoBody []byte, publicKey string) {
	var serversSshInfo []ServersSshInfo;
	if err := json.Unmarshal(serversInfoBody, &serversSshInfo); err != nil {
		fmt.Println(string(serversInfoBody))
		fmt.Println("Failed to parse response")
		os.Exit(exitCodeInvalidResponse)
	} else {
		sshPassword, _ := cmd.Flags().GetString("password")
		server := serversSshInfo[0].ExternalIp
		port := "22"
		server = server + ":" + port
		config := &ssh.ClientConfig{
			User: "root",
			Auth: []ssh.AuthMethod{
				ssh.Password(sshPassword),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
		conn, err := ssh.Dial("tcp", server, config)
		if err != nil {
			panic("Failed to dial: " + err.Error())
		}
		defer conn.Close()

		session, err := conn.NewSession()
		if err != nil {
			panic("Failed to create session: " + err.Error())
		}
		defer session.Close()

		session.Stdout = os.Stdout
		session.Stderr = os.Stderr

		if publicKey == "" {
			session.Stdin = os.Stdin

			modes := ssh.TerminalModes{
				ssh.ECHO: 0, // disable echoing
			}

			width, height, err := terminal.GetSize(int(os.Stdin.Fd()))

			if err != nil {
				fmt.Printf("Failed to initiate a terminal: %s\n", err)
				os.Exit(exitCodeUnexpected)
			}

			if err := session.RequestPty("xterm", width, height, modes); err != nil {
				fmt.Printf("request for pseudo terminal failed: %s", err)
				os.Exit(exitCodeUnexpected)
			}

			if err := session.Shell(); err != nil {
				fmt.Printf("failed to start shell: %s", err)
				os.Exit(exitCodeUnexpected)
			}

			if err = session.Wait(); err != nil {
				fmt.Printf("%s\n", err)
				os.Exit(exitCodeUnexpected)
			} else {
				os.Exit(0)
			}
		} else {
			publicKeyBytes, err := ioutil.ReadFile(publicKey)
			if err != nil {
				fmt.Printf("Failed to ready public key file\n")
				os.Exit(exitCodeUnexpected)
			}
			err = session.Run(fmt.Sprintf("echo '%s' >> ~/.ssh/authorized_keys", string(publicKeyBytes)))
			if err != nil {
				fmt.Printf("Failed to add public key to authorized keys\n")
				os.Exit(exitCodeUnexpected)
			}
			fmt.Printf("Successfuly added public key to authorized keys on the server\n")
			os.Exit(0)
		}
	}
}
