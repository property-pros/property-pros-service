package awss3

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

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
	bucket string
}

func NewClient() interfaces.IDocUploader {
	sess := session.Must(session.NewSession(&aws.Config{
		// TODO: move to config
		Endpoint:    aws.String("http://s3mock:9090"),
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials("accessKey", "secretKey", ""),
	}))

	cli := s3.New(sess, aws.NewConfig())

	return &AWSS3Client{
		client: cli,
		bucket: "documents",
	}
}

func (c *AWSS3Client) PutObject(ctx context.Context, content []byte) (string, error) {
	newKey := uuid.New().String()
	input := &s3.PutObjectInput{
		Bucket: aws.String(c.bucket),
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

func (c *AWSS3Client) GetObject(ctx context.Context, key string) ([]byte, error) {
	fmt.Printf("getting object for key: %v\n", key)
	input := &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}

	res, err := c.client.GetObject(input)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	fileContent, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading file content from S3 response:", err)
	}
	fmt.Println(string(fileContent))
	return fileContent, nil
}
