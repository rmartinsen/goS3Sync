package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	comparisons := GetSftpComparisons()
	localDir := "/Users/admin/go/src/github.com/rmartinsen/s3_sync/files/"

	timeStamp := time.Now().Format("20060102150405")

	for _, comparison := range comparisons {
		createLocalDir(localDir, comparison)
		existingObjectPath := filepath.Join(localDir, comparison.name, "existing")
		comparisonObjectPath := filepath.Join(localDir, comparison.name, "new")

		GetS3File(comparison, existingObjectPath)
		getSFTPFile(comparison, comparisonObjectPath)
		if compareFiles(existingObjectPath, comparisonObjectPath) == false {
			newFileName := comparison.s3Prefix + "_" + timeStamp
			UploadS3File(comparisonObjectPath, comparison.s3Bucket, newFileName)
		}
	}
}

func createLocalDir(localDir string, comparison SftpComparison) {
	path := filepath.Join(localDir, comparison.name)
	os.MkdirAll(path, os.ModePerm)
}

const chunkSize = 64000

func compareFiles(file1, file2 string) bool {
	// Check file size ...

	f1, err := os.Open(file1)
	if err != nil {
		log.Fatal(err)
	}

	f2, err := os.Open(file2)
	if err != nil {
		log.Fatal(err)
	}

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				log.Fatal(err1, err2)
			}
		}

		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}
