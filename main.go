package main

import (
	"os"
	"path/filepath"
)

func main() {
	comparisons := GetSftpComparisons()
	println(comparisons)
	localDir := "/Users/admin/go/src/github.com/rmartinsen/s3_sync/files/"
	for _, comparison := range comparisons {
		createLocalDir(localDir, comparison)
		getSFTPFile(comparison, localDir)
		GetS3File(comparison, localDir)
	}
}

func createLocalDir(localDir string, comparison SftpComparison) {
	path := filepath.Join(localDir, comparison.name)
	os.MkdirAll(path, os.ModePerm)
}
