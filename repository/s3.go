package repository

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/juju/errors"
)

const (
	S3_REGION = "us-east-2"
	S3_BUCKET = "cos143"
	S3_ACL    = "public-read"
)

type S3Response struct {
	FileURL string `json:"fileURL"`
}

type S3Repository struct {
	s3Session *s3.S3
}

func NewS3Repository(accessKey, secretKey string) (*S3Repository, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(S3_REGION),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		return nil, errors.Trace(err)
	}

	client := s3.New(sess)

	return &S3Repository{
		s3Session: client,
	}, nil
}

func (s *S3Repository) AddFileToS3(name string, reader *bytes.Reader) (string, error) {
	_, err := s.s3Session.PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(S3_BUCKET),
		Key:                aws.String(fmt.Sprintf("uploads/%s", name)),
		ACL:                aws.String(S3_ACL),
		CacheControl:       aws.String("private, max-age=31536000"),
		Body:               reader,
		ContentDisposition: aws.String("attachment"),
	})

	url := fmt.Sprintf("https://cos143.s3.us-east-2.amazonaws.com/uploads/%s", name)

	return url, errors.Trace(err)
}
