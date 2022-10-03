package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, s3Event events.S3Event) (string, error) {
	// The session the S3 Downloader will use
	sess := session.Must(session.NewSession())

	// Create a downloader with the session and default options
	downloader := s3manager.NewDownloader(sess)

	// Create a file to write the S3 Object contents to.
	filename := "/tmp/s3object"
	f, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("failed to create file %q, %v", filename, err)
	}
	defer f.Close()

	// Write the contents of S3 Object to the file
	record := s3Event.Records[0].S3
	n, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(record.Bucket.Name),
		Key:    aws.String(record.Object.Key),
	})
	if err != nil {
		return "", fmt.Errorf("failed to download file, %v", err)
	}
	log.Printf("file downloaded, %d bytes\n", n)

	//// Zip the file
	// Create a buffer to write our archive to.
	archivePath := fmt.Sprintf("/tmp/%s.zip", record.Object.Key)
	archive, err := os.Create(archivePath)
	if err != nil {
		log.Printf("Error creating archive: %s", err.Error())
		return "", fmt.Errorf("failed to create archive, %v", err)
	}
	defer archive.Close()

	// Create a new zip archive.
	w := zip.NewWriter(archive)

	// Add some files to the archive.
	w1, err := w.Create(record.Object.Key)
	if err != nil {
		log.Printf("Error adding file to archive: %s", err.Error())
		return "", fmt.Errorf("failed to add file to archive, %v", err)
	}

	_, err = io.Copy(w1, f)
	if err != nil {
		log.Printf("Error copying file to archive: %s", err.Error())
		return "", fmt.Errorf("failed to copying file to archive, %v", err)
	}

	// Make sure to check the error on Close.
	if err = w.Close(); err != nil {
		return "", fmt.Errorf("failed to copying file to archive, %v", err)
	}

	/// Upload zip to S3
	archiveToUpload, err := os.Open(archivePath)
	if err != nil {
		return "", fmt.Errorf("failed to open archive, %v", err)
	}

	defer archiveToUpload.Close()
	// Setup the S3 Upload Manager. Also see the SDK doc for the Upload Manager
	// for more information on configuring part size, and concurrency.
	//
	// http://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3manager/#NewUploader
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(record.Bucket.Name),
		Key:    aws.String(fmt.Sprintf("compressed/%s.zip", record.Object.Key)),
		Body:   archiveToUpload,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3, %v", err)
	}

	return fmt.Sprintf("Hello %s!", s3Event.Records[0].EventName), nil
}

func main() {
	lambda.Start(HandleRequest)
}
