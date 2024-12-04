package storage

import (
	"context"
	"io"
)

// Storage defines the interface for object storage operations
type Storage interface {
	// Core storage operations
	UploadModel(ctx context.Context, modelID string, data io.Reader, size int64) error
	DownloadModel(ctx context.Context, modelID string) (io.Reader, error)
	DeleteModel(ctx context.Context, modelID string) error

	// Optional: Add more methods as needed
	ListModels(ctx context.Context) ([]string, error)
	ModelExists(ctx context.Context, modelID string) (bool, error)
}

// type StorageConfig struct {
// 	Endpoint        string
// 	AccessKeyID     string
// 	SecretAccessKey string
// 	UseSSL          bool
// 	BucketName      string
// }

// func NewMinioClient(cfg StorageConfig) (*minio.Client, error) {
// 	// Initialize MinIO client
// 	client, err := minio.New(cfg.Endpoint, &minio.Options{
// 		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
// 		Secure: cfg.UseSSL,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create minio client: %v", err)
// 	}

// 	// Check if bucket exists and create it if it doesn't
// 	exists, err := client.BucketExists(context.Background(), cfg.BucketName)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to check bucket existence: %v", err)
// 	}

// 	if !exists {
// 		err = client.MakeBucket(context.Background(), cfg.BucketName, minio.MakeBucketOptions{})
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to create bucket: %v", err)
// 		}
// 	}

// 	return client, nil
// }
