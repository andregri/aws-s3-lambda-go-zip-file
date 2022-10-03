package main

import (
	"context"
	"fmt"
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
	return fmt.Sprintf("Hello %s!", s3Event.Records[0].EventName), nil
}

func main() {
	lambda.Start(HandleRequest)
}
