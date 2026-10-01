package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
	client            *s3.Client
	bucket, publicURL string
}

func New(ctx context.Context, endpoint, region, bucket, accessKey, secretKey, publicURL string) (*S3, error) {
	if bucket == "" {
		return nil, nil
	}
	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region), awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")))
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		cfg.BaseEndpoint = aws.String(endpoint)
	}
	client := s3.NewFromConfig(cfg, func(o *s3.Options) { o.UsePathStyle = endpoint != "" })
	if _, err := client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(bucket)}); err != nil {
		if _, createErr := client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(bucket)}); createErr != nil {
			return nil, fmt.Errorf("ensure S3 bucket %q: %w", bucket, createErr)
		}
	}
	return &S3{client: client, bucket: bucket, publicURL: strings.TrimRight(publicURL, "/")}, nil
}
func (s *S3) PutDoctorPhoto(ctx context.Context, key, contentType string, body io.Reader) (string, error) {
	if s == nil {
		return "", fmt.Errorf("S3 is not configured")
	}
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(key), Body: body, ContentType: aws.String(contentType)})
	if err != nil {
		return "", fmt.Errorf("put S3 object: %w", err)
	}
	if s.publicURL == "" {
		return key, nil
	}
	u, err := url.JoinPath(s.publicURL, path.Clean(key))
	return u, err
}
func (s *S3) Delete(ctx context.Context, key string) error {
	if s == nil {
		return fmt.Errorf("S3 is not configured")
	}
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(key)})
	return err
}
