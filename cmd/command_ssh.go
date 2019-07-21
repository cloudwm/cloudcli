package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

type ServersSshInfo struct {
	ExternalIp string `json:"externalIp"`
}

func commandRunSsh(cmd *cobra.Command, command SchemaCommand, serversInfoBody []byte) {
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
		session.Stdin = os.Stdin

		modes := ssh.TerminalModes{
			ssh.ECHO:          0,     // disable echoing
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
	}
}
