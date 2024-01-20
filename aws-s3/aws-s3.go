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
	"github.com/vireocloud/property-pros-service/config"
	"github.com/vireocloud/property-pros-service/interfaces"
)

type AWSS3Client struct {
	interfaces.IDocUploader
	client *s3.S3
	bucket string
}

// TODO: user config to connect to s3
func NewClient(config *config.Config) interfaces.IDocUploader {
	
	sess := session.Must(session.NewSession(&aws.Config{
		// TODO: move to config
		Endpoint: aws.String(config.S3Endpoint),
		Region:      aws.String(config.S3Region),
		Credentials: credentials.NewStaticCredentials(config.S3AccessKey, config.S3PrivateKey, ""),
		// needed for local docker image, should be false for aws s3
		S3ForcePathStyle: aws.Bool(config.S3ForcePathStyle),
	}))

	cli := s3.New(sess, aws.NewConfig())

	awsCli := &AWSS3Client{
		client: cli,
		bucket: config.S3BucketName,
	}

	err := awsCli.CreateBucketIfNotExists(context.TODO(), config.S3BucketName)
	if err != nil {
		panic(err)
	}

	return awsCli
}

func (c *AWSS3Client) CreateBucketIfNotExists(ctx context.Context, bucket string) error {
	fmt.Println("checking if the bucket exists")

	_, err := c.client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		fmt.Printf("creating bucket, as it doesn't exist; err: %v", err)
		_, err := c.client.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		})

		fmt.Printf("Error creating bucket: %v", err)

		// Wait for the bucket to be created before proceeding.
		err = c.client.WaitUntilBucketExists(&s3.HeadBucketInput{
			Bucket: aws.String(bucket),
		})

		if err != nil {
			fmt.Println("Error waiting for bucket:", err)
			return fmt.Errorf("Error waiting for bucket, %w", err)
		}
	}

	return nil
}

func (c *AWSS3Client) PutObject(ctx context.Context, content []byte) (string, error) {
	fmt.Println("putting object to s3")
	newKey := uuid.New().String()
	input := &s3.PutObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(newKey),
		Body:   bytes.NewReader(content),
	}

	_, err := c.client.PutObject(input)
	if err != nil {
		return "", err
	}

	return newKey, nil
}

func (c *AWSS3Client) GetObject(ctx context.Context, key string) ([]byte, error) {
	fmt.Println("getting object from s3")
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

	return fileContent, nil
}
