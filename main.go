package main

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/secrets", list_secrets_handler)
	http.HandleFunc("/addsecret", add_secret_form_handler)
	http.HandleFunc("/createsecret", add_secret_process_handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type SecretsPage struct {
	Secrets []SftpComparison
}

func list_secrets_handler(w http.ResponseWriter, r *http.Request) {

	secrets := GetSftpComparisons()

	p := &SecretsPage{Secrets: secrets}
	t, err := template.ParseFiles("viewsecrets.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, p)
}

func add_secret_form_handler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("addsecret.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, nil)
}

func add_secret_process_handler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("secretadded.html")
	if err != nil {
		log.Fatal(err)
	}

	err = r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	form := r.PostForm

	port, err := strconv.Atoi(form["port"][0])
	if err != nil {
		log.Fatal(err)
	}

	comparison := SftpComparison{
		Name:     form["secretname"][0],
		Host:     form["sftphost"][0],
		Username: form["username"][0],
		Password: form["password"][0],
		SftpPath: form["sftppath"][0],
		S3Bucket: form["s3bucket"][0],
		S3Prefix: form["s3prefix"][0],
		Port:     port,
	}

	CreateSFTPComparison(comparison)
	t.Execute(w, nil)
}

func DoComparisons() {
	comparisons := GetSftpComparisons()
	localDir := "/Users/admin/go/src/github.com/rmartinsen/s3_sync/files/"

	timeStamp := time.Now().Format("20060102150405")

	for _, comparison := range comparisons {
		createLocalDir(localDir, comparison)
		existingObjectPath := filepath.Join(localDir, comparison.Name, "existing")
		comparisonObjectPath := filepath.Join(localDir, comparison.Name, "new")

		GetS3File(comparison, existingObjectPath)
		getSFTPFile(comparison, comparisonObjectPath)
		if compareFiles(existingObjectPath, comparisonObjectPath) == false {
			newFileName := comparison.S3Prefix + "_" + timeStamp
			UploadS3File(comparisonObjectPath, comparison.S3Bucket, newFileName)
		}
	}
}

func createLocalDir(localDir string, comparison SftpComparison) {
	path := filepath.Join(localDir, comparison.Name)
	os.MkdirAll(path, os.ModePerm)
}

const chunkSize = 64000

func compareFiles(file1, file2 string) bool {

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
