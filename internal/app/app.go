package app

import (
	"context"
	"io"
	"naivebayesservice/internal/config"
	"naivebayesservice/internal/storage"
)

type App struct {
	storage storage.Storage
	// Add other dependencies here (e.g., model trainer, API handlers)
}

func New(cfg *config.Config) (*App, error) {
	// Initialize storage
	storageClient, err := storage.NewMinioClient(storage.StorageConfig{
		Endpoint:        cfg.Storage.Endpoint,
		AccessKeyID:     cfg.Storage.AccessKeyID,
		SecretAccessKey: cfg.Storage.SecretAccessKey,
		UseSSL:          cfg.Storage.UseSSL,
		BucketName:      cfg.Storage.BucketName,
	})
	if err != nil {
		return nil, err
	}

	return &App{
		storage: storageClient,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	// Start your application (e.g., HTTP server)
	return nil
}

// Add these methods to expose storage operations
func (a *App) UploadModel(ctx context.Context, modelID string, data io.Reader, size int64) error {
	return a.storage.UploadModel(ctx, modelID, data, size)
}

func (a *App) DownloadModel(ctx context.Context, modelID string) (io.Reader, error) {
	return a.storage.DownloadModel(ctx, modelID)
}

func (a *App) DeleteModel(ctx context.Context, modelID string) error {
	return a.storage.DeleteModel(ctx, modelID)
}

func (a *App) ListModels(ctx context.Context) ([]string, error) {
	return a.storage.ListModels(ctx)
}

func (a *App) ModelExists(ctx context.Context, modelID string) (bool, error) {
	return a.storage.ModelExists(ctx, modelID)
}
