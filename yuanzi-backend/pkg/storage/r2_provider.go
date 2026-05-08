package storage

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	appConfig "yuanzi-backend/config"
)

// R2Provider implements Provider using Cloudflare R2 (S3-compatible).
type R2Provider struct {
	client    *s3.Client
	presigner *s3.PresignClient
	bucket    string
	publicURL string
	accountID string
}

// NewR2Provider creates a Provider backed by Cloudflare R2.
// Configuration is read from the global config (config.GlobalConfig.R2) first,
// then falls back to environment variables for backward compatibility.
//
// Expected env vars (fallback):
//
//	R2_ACCOUNT_ID, R2_ACCESS_KEY_ID, R2_ACCESS_KEY_SECRET, R2_BUCKET, R2_PUBLIC_URL.
func NewR2Provider() (Provider, error) {
	accountID := appConfig.GlobalConfig.R2.AccountID
	accessKeyID := appConfig.GlobalConfig.R2.AccessKeyID
	accessKeySecret := appConfig.GlobalConfig.R2.AccessKeySecret
	bucket := appConfig.GlobalConfig.R2.Bucket
	publicURL := appConfig.GlobalConfig.R2.PublicURL

	// Fallback to environment variables if config values are empty
	if accountID == "" {
		accountID = os.Getenv("R2_ACCOUNT_ID")
	}
	if accessKeyID == "" {
		accessKeyID = os.Getenv("R2_ACCESS_KEY_ID")
	}
	if accessKeySecret == "" {
		accessKeySecret = os.Getenv("R2_ACCESS_KEY_SECRET")
	}
	if bucket == "" {
		bucket = os.Getenv("R2_BUCKET")
	}
	if publicURL == "" {
		publicURL = os.Getenv("R2_PUBLIC_URL")
	}
	if publicURL == "" {
		publicURL = os.Getenv("R2_CUSTOM_DOMAIN")
	}

	if accountID == "" || accessKeyID == "" || accessKeySecret == "" || bucket == "" {
		return nil, fmt.Errorf("missing required R2 configuration: account_id, access_key_id, access_key_secret, bucket")
	}

	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("auto"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKeyID, accessKeySecret, ""),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load R2 config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	presigner := s3.NewPresignClient(client)

	return &R2Provider{
		client:    client,
		presigner: presigner,
		bucket:    bucket,
		publicURL: publicURL,
		accountID: accountID,
	}, nil
}

// GetUploadSignature generates a presigned PUT URL for direct upload to R2.
func (p *R2Provider) GetUploadSignature(key string, maxSize int64, expireSeconds int) (*UploadSignature, error) {
	if expireSeconds <= 0 {
		expireSeconds = 300
	}

	ctx := context.Background()
	req, err := p.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(time.Duration(expireSeconds)*time.Second))
	if err != nil {
		return nil, fmt.Errorf("failed to presign PUT URL: %w", err)
	}

	expiresAt := time.Now().Add(time.Duration(expireSeconds) * time.Second).Unix()

	return &UploadSignature{
		PresignedURL: req.URL,
		UploadURL:    req.URL,
		AccessURL:    p.GetURL(key),
		ThumbnailURL: p.GetThumbnailURL(key, 0),
		ExpiresAt:    expiresAt,
	}, nil
}

// GetURL returns the public access URL for an R2 object.
func (p *R2Provider) GetURL(key string) string {
	if p.publicURL != "" {
		return fmt.Sprintf("%s/%s", p.publicURL, key)
	}
	return fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s/%s", p.accountID, p.bucket, key)
}

// GetThumbnailURL returns a URL that may trigger R2 Image Resizing.
// If width > 0, appends an image-resizing query parameter.
func (p *R2Provider) GetThumbnailURL(key string, width int) string {
	url := p.GetURL(key)
	if width > 0 {
		return url + fmt.Sprintf("?x-oss-process=image/resize,w_%d", width)
	}
	return url
}

// DeleteObject removes the object from R2 via the S3 DeleteObject API.
func (p *R2Provider) DeleteObject(key string) error {
	if key == "" {
		return fmt.Errorf("object key cannot be empty")
	}
	ctx := context.Background()
	_, err := p.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete R2 object: %w", err)
	}
	return nil
}
