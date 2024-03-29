package common

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"os"
	"strings"
)

func CheckFileOverSSH(filePath string, user string, host string, privateKeyFilePath string, hostKeyFilePath string) bool {
	key, err := os.ReadFile(privateKeyFilePath)
	if err != nil {
		log.WithFields(log.Fields{
			"component": "ssh",
			"event":     "auth",
		}).Fatalf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.WithFields(log.Fields{
			"component": "ssh",
			"event":     "auth",
		}).Fatalf("unable to parse private key: %v", err)
	}

	hostKeyCallback, err := knownhosts.New(hostKeyFilePath)
	if err != nil {
		log.WithFields(log.Fields{
			"component": "ssh",
			"event":     "auth",
		}).Fatal(err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,
	}

	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		log.WithFields(log.Fields{
			"component": "ssh",
			"event":     "connect",
		}).Fatalf("unable to connect: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.WithFields(log.Fields{
			"component": "ssh",
			"event":     "connect",
		}).Fatalf("unable to create session: %s", err)
	}
	defer session.Close()

	b, err := session.Output(
		fmt.Sprintf("[[ -f %s ]] && echo \"File exists\" || echo \"File does not exist\"", filePath))
	if err != nil {
		log.WithFields(log.Fields{
			"component": "ssh",
			"event":     "remote-run",
		}).Fatalf("failed to execute: %s", err)
	}
	outputString := strings.TrimSuffix(string(b), "\n")

	if outputString == "File exists" {
		return true
	}
	return false
}
