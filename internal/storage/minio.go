package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioStorage implements the Storage interface
type MinioStorage struct {
	client     *minio.Client
	bucketName string
}

func NewMinioClient(cfg StorageConfig) (*MinioStorage, error) {
	// Initialize MinIO client
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %v", err)
	}

	// Check if bucket exists and create it if it doesn't
	exists, err := client.BucketExists(context.Background(), cfg.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %v", err)
	}

	if !exists {
		err = client.MakeBucket(context.Background(), cfg.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %v", err)
		}
	}

	return &MinioStorage{
		client:     client,
		bucketName: cfg.BucketName,
	}, nil
}

// Implement Storage interface methods
func (m *MinioStorage) UploadModel(ctx context.Context, modelID string, data io.Reader, size int64) error {
	_, err := m.client.PutObject(ctx, m.bucketName, modelID, data, size, minio.PutObjectOptions{})
	return err
}

func (m *MinioStorage) DownloadModel(ctx context.Context, modelID string) (io.Reader, error) {
	obj, err := m.client.GetObject(ctx, m.bucketName, modelID, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get model: %v", err)
	}
	return obj, nil
}

func (m *MinioStorage) DeleteModel(ctx context.Context, modelID string) error {
	return m.client.RemoveObject(ctx, m.bucketName, modelID, minio.RemoveObjectOptions{})
}

func (m *MinioStorage) ListModels(ctx context.Context) ([]string, error) {
	var models []string
	objects := m.client.ListObjects(ctx, m.bucketName, minio.ListObjectsOptions{})
	for object := range objects {
		if object.Err != nil {
			return nil, object.Err
		}
		models = append(models, object.Key)
	}
	return models, nil
}

func (m *MinioStorage) ModelExists(ctx context.Context, modelID string) (bool, error) {
	_, err := m.client.StatObject(ctx, m.bucketName, modelID, minio.StatObjectOptions{})
	if err != nil {
		// Check if the error indicates that the object doesn't exist
		if errResponse := minio.ToErrorResponse(err); errResponse.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
