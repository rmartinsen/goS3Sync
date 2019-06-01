package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func getMostRecentKey(bucket string, prefix string) string {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	svc := s3.New(sess)

	result, err := svc.ListBuckets(nil)

	if err != nil {
		log.Fatal(err)
	}

	println(result)
	return "doneskis"
}

func main() {
	getMostRecentKey("1", "2")
}
