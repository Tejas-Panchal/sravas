package service

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Storage stores video files in an S3 bucket
type S3Storage struct {
	client *s3.Client
	bucket string
}

// NewS3Storage creates an S3Storage connected to the given bucket and region
func NewS3Storage(bucket, region string) (*S3Storage, error) {
	if bucket == "" {
		return nil, fmt.Errorf("S3_BUCKET is required")
	}
	if region == "" {
		region = "us-east-1"
	}

	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &S3Storage{
		client: s3.NewFromConfig(cfg),
		bucket: bucket,
	}, nil
}

// Save uploads a file to the S3 bucket
func (s *S3Storage) Save(key string, reader io.Reader) error {
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   reader,
	})
	return err
}

// Get retrieves a file from the S3 bucket
func (s *S3Storage) Get(key string) (io.ReadCloser, error) {
	out, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}

// Delete removes a file from the S3 bucket
func (s *S3Storage) Delete(key string) error {
	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}
