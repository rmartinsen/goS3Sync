package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type SftpComparison struct {
	Name     string
	Host     string
	Username string
	Password string
	SftpPath string
	S3Bucket string
	S3Prefix string
	Port     int
}

func GetSftpComparisons() []SftpComparison {

	parameterOutput := getSecrets()
	comparisons := make([]SftpComparison, 0)

	for _, parameter := range parameterOutput.Parameters {
		comparison := &SftpComparison{}
		quotedParam := "`" + string(*parameter.Value) + "`"
		param, err := strconv.Unquote(quotedParam)
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal([]byte(param), comparison)
		if err != nil {
			log.Fatal(err)
		}

		comparisons = append(comparisons, *comparison)
	}
	// sameCompare := SftpComparison{
	// 	"same_compare", "localhost", "admin", "admin",
	// 	"/Users/admin/Documents/test_1",
	// 	"compare-test-1", "test", 22}

	// diffCompare := SftpComparison{
	// 	"diff_compare", "localhost", "admin", "admin",
	// 	"/Users/admin/Documents/test_2",
	// 	"compare-test-2", "test", 22}

	// results = append(results, sameCompare)
	// results = append(results, diffCompare)

	return comparisons
}

func getSecrets() ssm.GetParametersByPathOutput {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	svc := ssm.New(sess)

	t := true

	results, _ := svc.GetParametersByPath(
		&ssm.GetParametersByPathInput{
			Path:           aws.String("/file_compare/"),
			WithDecryption: &t,
		})

	return *results
}

func CreateSFTPComparison(comparison SftpComparison) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	svc := ssm.New(sess)

	comparisonName := fmt.Sprintf("/file_compare/%s", comparison.Name)
	comparisonValue := jsonify(comparison)

	svc.PutParameter(&ssm.PutParameterInput{
		Name:  aws.String(comparisonName),
		Value: aws.String(comparisonValue),
		Type:  aws.String("SecureString"),
	})

}

func jsonify(comparison SftpComparison) string {
	return fmt.Sprintf(`
	{
		"name": "%s",
		"host": "%s",
		"username": "%s",
		"password": "%s",
		"sftpPath": "%s",
		"s3Bucket": "%s",
		"s3Prefix": "%s",
		"port": %d
	}
	`, comparison.Name, comparison.Host, comparison.Username, comparison.Password,
		comparison.SftpPath, comparison.S3Bucket, comparison.S3Prefix, comparison.Port)

}
