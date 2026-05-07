// Package storage provides an abstraction layer for object storage backends,
// supporting both Alibaba Cloud OSS and Cloudflare R2.
package storage

import (
	"fmt"
	"yuanzi-backend/config"
)

// NewProviderFromConfig returns the appropriate Provider based on the global config.
// It defaults to OSS if the provider setting is unrecognized.
func NewProviderFromConfig() (Provider, error) {
	provider := config.GlobalConfig.Storage.Provider
	if provider == "" {
		provider = "oss"
	}

	switch provider {
	case "r2":
		return NewR2Provider()
	case "oss":
		return NewOSSProvider(), nil
	default:
		return nil, fmt.Errorf("unsupported storage provider: %s", provider)
	}
}

// UploadSignature holds the information needed for a client to upload an object
// directly to the storage backend.
type UploadSignature struct {
	FormData     map[string]string // for OSS multipart POST
	PresignedURL string            // for R2 presigned PUT
	UploadURL    string            // upload endpoint URL
	AccessURL    string            // final access URL
	ThumbnailURL string            // thumbnail access URL
	ExpiresAt    int64             // unix timestamp
}

// Provider defines the common interface for storage backends.
type Provider interface {
	// GetUploadSignature generates upload credentials or a presigned URL for the given object key.
	GetUploadSignature(key string, maxSize int64, expireSeconds int) (*UploadSignature, error)
	// GetURL returns the public access URL for an object.
	GetURL(key string) string
	// GetThumbnailURL returns the URL for a resized thumbnail of the object.
	GetThumbnailURL(key string, width int) string
	// DeleteObject removes the object from storage.
	DeleteObject(key string) error
}
