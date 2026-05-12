package storage

import (
	"time"

	"yuanzi-backend/pkg/oss"
)

// OSSProvider implements Provider using the legacy Alibaba Cloud OSS client.
type OSSProvider struct {
	client *oss.Client
}

// NewOSSProvider creates a Provider backed by Alibaba Cloud OSS.
func NewOSSProvider() Provider {
	return &OSSProvider{
		client: oss.NewClient(),
	}
}

// GetUploadSignature generates an OSS Post signature for browser-based multipart uploads.
func (p *OSSProvider) GetUploadSignature(key string, maxSize int64, expireSeconds int, options ...UploadOption) (*UploadSignature, error) {
	if expireSeconds <= 0 {
		expireSeconds = 300
	}

	opts := applyUploadOptions(options)
	sig, uploadURL, err := p.client.GetPostSignature(key, maxSize, expireSeconds)
	if err != nil {
		return nil, err
	}

	formData := map[string]string{
		"OSSAccessKeyId": sig.OSSAccessKeyID,
		"policy":         sig.Policy,
		"signature":      sig.Signature,
		"key":            key,
	}
	if opts.ContentType != "" {
		formData["Content-Type"] = opts.ContentType
	}

	expiresAt := time.Now().Add(time.Duration(expireSeconds) * time.Second).Unix()

	return &UploadSignature{
		FormData:     formData,
		UploadURL:    uploadURL,
		AccessURL:    p.client.GetURL(key),
		ThumbnailURL: p.client.GetThumbnailURL(key, 0),
		ExpiresAt:    expiresAt,
	}, nil
}

// GetURL returns the public URL for the given OSS object key.
func (p *OSSProvider) GetURL(key string) string {
	return p.client.GetURL(key)
}

// GetThumbnailURL returns a URL with OSS image-processing parameters for a thumbnail.
func (p *OSSProvider) GetThumbnailURL(key string, width int) string {
	return p.client.GetThumbnailURL(key, width)
}

// DeleteObject removes the object from OSS.
func (p *OSSProvider) DeleteObject(key string) error {
	return p.client.DeleteObject(key)
}
