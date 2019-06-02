package main

import (
	"log"
	"os"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func getSFTPFile(comparison SftpComparison, localDir string) {
	config := &ssh.ClientConfig{
		User: comparison.username,
		Auth: []ssh.AuthMethod{
			ssh.Password(comparison.password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	ssh, err := ssh.Dial("tcp", "127.0.0.1:22", config)

	if err != nil {
		log.Fatal(err)
	}

	sftp, err := sftp.NewClient(ssh)

	if err != nil {
		log.Fatal(err)
	}
	defer sftp.Close()

	source, err := sftp.Open(comparison.sftpPath)
	if err != nil {
		log.Fatal(err)
	}
	defer source.Close()

	splitPath := strings.Split(comparison.sftpPath, "/")

	sftpFileName := splitPath[len(splitPath)-1]

	target, err := os.Create(localDir + comparison.name + "/" + "sftp_" + sftpFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer target.Close()

	source.WriteTo(target)
}
