// Package s3 provides the admin module's gateway to Yandex Object Storage.
package s3

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/nagorneva/nagorneva_reviews/internal/config"
)

const defaultContentType = "application/octet-stream"

// client contains only the S3 operations used by the admin module.
type client interface {
	PutObject(context.Context, *awss3.PutObjectInput, ...func(*awss3.Options)) (*awss3.PutObjectOutput, error)
	DeleteObject(context.Context, *awss3.DeleteObjectInput, ...func(*awss3.Options)) (*awss3.DeleteObjectOutput, error)
}

// Gateway hides S3 client construction and object-storage calls from actions.
type Gateway struct {
	client            client
	bucket, publicURL string
}

// NewGateway creates the admin module's S3 gateway from application config.
// A nil gateway means S3 is disabled for this module.
func NewGateway(ctx context.Context, cfg config.Config) (*Gateway, error) {
	if cfg.S3Bucket == "" {
		return nil, nil
	}
	if cfg.S3AccessKey == "" || cfg.S3SecretKey == "" {
		return nil, fmt.Errorf("S3_ACCESS_KEY and S3_SECRET_KEY are required when S3_BUCKET is set")
	}

	options := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(cfg.S3Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.S3AccessKey, cfg.S3SecretKey, "")),
	}
	if cfg.S3Endpoint != "" {
		options = append(options, awsconfig.WithEndpointResolver(aws.EndpointResolverFunc(
			func(_, _ string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: cfg.S3Endpoint, SigningRegion: cfg.S3Region}, nil
			},
		)))
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("load S3 configuration: %w", err)
	}

	publicURL := cfg.S3PublicURL
	if publicURL == "" {
		publicURL = "https://storage.yandexcloud.net/" + cfg.S3Bucket
	}
	return &Gateway{
		client:    awss3.NewFromConfig(awsCfg),
		bucket:    cfg.S3Bucket,
		publicURL: strings.TrimRight(publicURL, "/"),
	}, nil
}

// PutDoctorPhoto uploads a doctor photo and returns its public URL.
func (g *Gateway) PutDoctorPhoto(ctx context.Context, key, filename string, body io.Reader) (string, error) {
	if g == nil {
		return "", fmt.Errorf("S3 is not configured")
	}
	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType == "" {
		contentType = defaultContentType
	}
	_, err := g.client.PutObject(ctx, &awss3.PutObjectInput{
		Bucket:      aws.String(g.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
		Metadata: map[string]string{
			"uploaded_at": time.Now().UTC().Format(time.RFC3339),
			"origin_name": filename,
		},
	})
	if err != nil {
		return "", fmt.Errorf("upload doctor photo: %w", err)
	}
	return url.JoinPath(g.publicURL, key)
}

// Delete removes an object from the module's S3 bucket.
func (g *Gateway) Delete(ctx context.Context, key string) error {
	if g == nil {
		return fmt.Errorf("S3 is not configured")
	}
	_, err := g.client.DeleteObject(ctx, &awss3.DeleteObjectInput{Bucket: aws.String(g.bucket), Key: aws.String(key)})
	if err != nil {
		return fmt.Errorf("delete S3 object: %w", err)
	}
	return nil
}
