package imagedata

import (
	"bytes"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/imgproxy/imgproxy/v3/config"
)

var aswClient *s3.S3

func newAWSConnection() error {
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(
			config.AWS.AccessKeyID,
			config.AWS.SecretAccessKey,
			""),
		Endpoint: aws.String("https://ams3.digitaloceanspaces.com"),
		Region:   aws.String("us-east-1"),
	}
	newSession := session.New(s3Config)
	aswClient = s3.New(newSession)

	_, err := aswClient.ListBuckets(nil) // ping for success connect
	if err != nil {
		return err
	}
	return nil
}

func UploadImage(awsClient *s3.S3, path string, img *ImageData) error {
	object := s3.PutObjectInput{
		Bucket:             aws.String(config.AWS.Bucket),
		Key:                aws.String(path),
		Body:               bytes.NewReader(img.Data),
		ContentType:        aws.String(http.DetectContentType(img.Data)),
		ContentDisposition: aws.String("attachment"),
		ACL:                aws.String("public-read"),
	}
	_, err := awsClient.PutObject(&object)
	if err != nil {
		return err
	}
	return nil
}
