package main

type SftpComparison struct {
	name     string
	host     string
	username string
	password string
	sftpPath string
	s3Bucket string
	s3Prefix string
	port     int
}

func GetSftpComparisons() []SftpComparison {

	results := make([]SftpComparison, 0)

	sameCompare := SftpComparison{
		"same_compare", "localhost", "admin", "admin",
		"/Users/admin/Documents/test_1",
		"compare-test-1", "test", 22}

	diffCompare := SftpComparison{
		"diff_compare", "localhost", "admin", "admin",
		"/Users/admin/Documents/test_2",
		"compare-test-2", "test", 22}

	results = append(results, sameCompare)
	results = append(results, diffCompare)

	return results
}
