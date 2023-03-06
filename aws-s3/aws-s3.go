package awss3

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/vireocloud/property-pros-service/interfaces"
)

type AWSS3Client struct {
	interfaces.IDocUploader
	client *s3.S3
}

func NewClient() interfaces.IDocUploader {
	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint:    aws.String("http://localhost:9090"),
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials("accessKey", "secretKey", ""),
	}))

	cli := s3.New(sess, aws.NewConfig())

	return &AWSS3Client{
		client: cli,
	}
}

func (c *AWSS3Client) PutObject(ctx context.Context, content []byte) (string, error) {
	fmt.Println("here in put object")
	newKey := uuid.New().String()
	input := &s3.PutObjectInput{
		Bucket: aws.String("documents"),
		Key:    aws.String(newKey),
		Body:   bytes.NewReader(content),
	}

	res, err := c.client.PutObject(input)
	if err != nil {
		return "", err
	}

	fmt.Printf("response from s3 is: %v \n", res)
	return newKey, nil
}
