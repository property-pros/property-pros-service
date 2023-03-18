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

var AWS_BUCKET = "documents"

func NewClient() interfaces.IDocUploader {
	sess := session.Must(session.NewSession(&aws.Config{
		// TODO: move to config
		Endpoint:    aws.String("http://s3mock:9090"),
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials("accessKey", "secretKey", ""),
		// needed for local docker image, should work for aws s3 too
		S3ForcePathStyle: aws.Bool(true),
	}))

	cli := s3.New(sess, aws.NewConfig())

	awsCli := &AWSS3Client{
		client: cli,
		bucket: AWS_BUCKET,
	}

	err := awsCli.CreateBucketIfNotExists(context.TODO(), AWS_BUCKET)
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
		fmt.Println("creating bucket, as it doesn't exist")
		_, err := c.client.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		})

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
