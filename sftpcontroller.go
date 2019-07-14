package main

import (
	"log"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func getSFTPFile(comparison SftpComparison, localPath string) {
	config := &ssh.ClientConfig{
		User: comparison.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(comparison.Password),
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

	source, err := sftp.Open(comparison.SftpPath)
	if err != nil {
		log.Fatal(err)
	}
	defer source.Close()

	target, err := os.Create(localPath)
	if err != nil {
		log.Fatal(err)
	}
	defer target.Close()

	source.WriteTo(target)
}
