package goS3

import (
	"bytes"
	"context"
	"fmt"
	"mime"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Client struct {
	client     *s3.Client
	presign    *s3.PresignClient
	bucketName string
}

type Config struct {
	AccessKey  string
	SecretKey  string
	Region     string
	BucketName string
}

func New(ctx context.Context, cfg Config) (*Client, error) {
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AccessKey,
				cfg.SecretKey,
				"",
			),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	s3Client := s3.NewFromConfig(awsCfg)

	return &Client{
		client:     s3Client,
		presign:    s3.NewPresignClient(s3Client),
		bucketName: cfg.BucketName,
	}, nil
}

func (c *Client) Upload(ctx context.Context, key string, file []byte) error {
	_, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:               aws.String(c.bucketName),
		Key:                  aws.String(key),
		Body:                 bytes.NewReader(file),
		ContentType:          aws.String(getContentType(key)),
		ContentDisposition:   aws.String("inline"),
		ServerSideEncryption: types.ServerSideEncryptionAes256,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	return nil
}

func (c *Client) GeneratePresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	req, err := c.presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiry
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return req.URL, nil
}

func getContentType(filename string) string {
	ext := filepath.Ext(filename)
	if contentType := mime.TypeByExtension(ext); contentType != "" {
		return contentType
	}
	return "application/octet-stream"
}
