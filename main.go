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
		existingObjectPath := filepath.Join(localDir, comparison.name, "existing")
		comparisonObjectPath := filepath.Join(localDir, comparison.name, "new")

		GetS3File(comparison, existingObjectPath)
		getSFTPFile(comparison, comparisonObjectPath)
		if CompareFiles(existingObjectPath, comparisonObjectPath) == false {
			newFileName := "NewTestOne"
			UploadS3File(existingObjectPath, comparison.s3Bucket, newFileName)
		}
	}
}

func createLocalDir(localDir string, comparison SftpComparison) {
	path := filepath.Join(localDir, comparison.name)
	os.MkdirAll(path, os.ModePerm)
}

func CompareFiles(firstFilePath string, secondFilePath string) bool {
	return false
}
