package main

import (
	"log"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func main() {
	comparisons := GetSftpComparisons()
	println(comparisons)
	for _, comparison := range comparisons {
		getSFTPFile(comparison)
	}
}

func getSFTPFile(comparison SftpComparison) {
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

	target, err := os.Create("/Users/admin/go/src/github.com/rmartinsen/s3_sync/files/" + comparison.name)
	if err != nil {
		log.Fatal(err)
	}
	defer target.Close()

	source.WriteTo(target)
}
