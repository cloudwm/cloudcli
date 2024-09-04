package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/moby/term"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
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

func startSSHSession(session *ssh.Session) error {
	var width, height int
	if runtime.GOOS == "windows" {
		width, height = 80, 40
	} else {
		winsize, err := term.GetWinsize(os.Stdin.Fd())
		if err != nil {
			return fmt.Errorf("failed to get terminal size: %v", err)
		}
		width, height = int(winsize.Width), int(winsize.Height)
	}
	modes := ssh.TerminalModes{ssh.ECHO: 0}
	if err := session.RequestPty("xterm", width, height, modes); err != nil {
		return fmt.Errorf("request for pseudo terminal failed: %v", err)
	}
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin
	if err := session.Shell(); err != nil {
		return fmt.Errorf("failed to start shell: %v", err)
	}
	return session.Wait()
}

func commandRunSsh(cmd *cobra.Command, command SchemaCommand, serversInfoBody []byte, publicKey string) {
	sshPassword, _ := cmd.Flags().GetString("password")
	sshPrivateKey, _ := cmd.Flags().GetString("key")
	if sshPassword != "" && sshPrivateKey != "" {
		fmt.Println("Must use either --password or --key, but not both")
		os.Exit(exitCodeInvalidFlags)
	} else if sshPassword == "" && sshPrivateKey == "" {
		fmt.Println("Must set either --password or --key")
		os.Exit(exitCodeInvalidFlags)
	}
	var serversSshInfo []ServersSshInfo
	if err := json.Unmarshal(serversInfoBody, &serversSshInfo); err != nil {
		fmt.Println(string(serversInfoBody))
		fmt.Println("Failed to parse response")
		os.Exit(exitCodeInvalidResponse)
	} else {
		var config ssh.ClientConfig
		if sshPrivateKey != "" {
			pkBytes, err := ioutil.ReadFile(sshPrivateKey)
			if err != nil {
				fmt.Printf("Failed to read private key file: %s\n", err.Error())
				os.Exit(exitCodeUnexpected)
			}
			signer, err := ssh.ParsePrivateKey(pkBytes)
			if err != nil {
				fmt.Printf("Unable to parse private key: %s\n", err.Error())
				os.Exit(exitCodeUnexpected)
			}
			config = ssh.ClientConfig{
				User: "root",
				Auth: []ssh.AuthMethod{
					ssh.PublicKeys(signer),
				},
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			}
		} else {
			config = ssh.ClientConfig{
				User: "root",
				Auth: []ssh.AuthMethod{
					ssh.Password(sshPassword),
				},
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			}
		}
		server := serversSshInfo[0].ExternalIp
		port := "22"
		server = server + ":" + port
		conn, err := ssh.Dial("tcp", server, &config)
		if err != nil {
			fmt.Printf("Failed to connect to the server: %s\n", err.Error())
			os.Exit(exitCodeInvalidResponse)
		}
		defer conn.Close()

		session, err := conn.NewSession()
		if err != nil {
			fmt.Printf("Failed to initiate SSH session: %s\n", err.Error())
			os.Exit(exitCodeInvalidResponse)
		}
		defer session.Close()

		if publicKey == "" {
			if err = startSSHSession(session); err != nil {
				fmt.Printf("%s\n", err)
				os.Exit(exitCodeUnexpected)
			} else {
				os.Exit(0)
			}
		} else {
			session.Stdout = os.Stdout
			session.Stderr = os.Stderr
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
