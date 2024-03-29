package main

import (
	"bytes"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func GetS3File(comparison SftpComparison, localPath string) {
	bucket := comparison.S3Bucket
	prefix := comparison.S3Prefix

	key := getMostRecentKey(bucket, prefix)

	file, err := os.Create(localPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	downloader := s3manager.NewDownloader(sess)

	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	println("Downloaded", bucket, key, "to", "localPath")

}

func getMostRecentKey(bucket string, prefix string) string {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	svc := s3.New(sess)

	bucketInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}

	result, err := svc.ListObjectsV2(bucketInput)

	if err != nil {
		log.Fatal(err)
	}

	maxKey := ""
	for _, result := range result.Contents {
		if *result.Key > maxKey {
			maxKey = *result.Key
		}
	}

	return maxKey
}

func UploadS3File(localPath string, bucket string, key string) {
	println("Uploading localpath to " + bucket + "/" + key)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(localPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)

	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        fileBytes,
		ContentType: aws.String(fileType),
	})

	if err != nil {
		log.Fatal(err)
	}
	println("File uploaded from " + localPath + "to: " + bucket + key)
}
