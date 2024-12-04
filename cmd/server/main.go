package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"naivebayesservice/internal/app"
	"naivebayesservice/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create application
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown gracefully
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	// Test storage operations
	if err := testStorage(ctx, application); err != nil {
		log.Fatalf("Storage test failed: %v", err)
	}
}

func testStorage(ctx context.Context, app *app.App) error {
	// We need to add a method to access storage from App
	// Let's update app.go first to add this method

	// Test uploading a model
	testData := []byte("test model data")
	modelID := fmt.Sprintf("test-model-%d", time.Now().Unix())

	fmt.Printf("Testing model upload with ID: %s\n", modelID)

	if err := app.UploadModel(ctx, modelID, bytes.NewReader(testData), int64(len(testData))); err != nil {
		return fmt.Errorf("failed to upload model: %v", err)
	}
	fmt.Println("Upload successful")

	// Test model exists
	exists, err := app.ModelExists(ctx, modelID)
	if err != nil {
		return fmt.Errorf("failed to check model existence: %v", err)
	}
	fmt.Printf("Model exists: %v\n", exists)

	// Test downloading the model
	reader, err := app.DownloadModel(ctx, modelID)
	if err != nil {
		return fmt.Errorf("failed to download model: %v", err)
	}

	downloadedData := new(bytes.Buffer)
	_, err = downloadedData.ReadFrom(reader)
	if err != nil {
		return fmt.Errorf("failed to read downloaded data: %v", err)
	}

	if string(downloadedData.Bytes()) != string(testData) {
		return fmt.Errorf("downloaded data doesn't match uploaded data")
	}
	fmt.Println("Download successful and data matches")

	// Test listing models
	models, err := app.ListModels(ctx)
	if err != nil {
		return fmt.Errorf("failed to list models: %v", err)
	}
	fmt.Printf("Models in storage: %v\n", models)

	// Test deleting the model
	if err := app.DeleteModel(ctx, modelID); err != nil {
		return fmt.Errorf("failed to delete model: %v", err)
	}
	fmt.Println("Delete successful")

	// Verify deletion
	exists, err = app.ModelExists(ctx, modelID)
	if err != nil {
		return fmt.Errorf("failed to check model existence after deletion: %v", err)
	}
	if exists {
		return fmt.Errorf("model still exists after deletion")
	}
	fmt.Println("Deletion verified")

	return nil
}

// package main

// import (
// 	"context"
// 	"log"
// 	"os"
// 	"os/signal"
// 	"syscall"

// 	"naivebayesservice/internal/app"
// 	"naivebayesservice/internal/config"
// )

// func main() {
// 	// Load configuration
// 	cfg, err := config.Load()
// 	if err != nil {
// 		log.Fatalf("Failed to load config: %v", err)
// 	}

// 	// Create application
// 	application, err := app.New(cfg)
// 	if err != nil {
// 		log.Fatalf("Failed to create application: %v", err)
// 	}

// 	// Setup context with cancellation
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	// Handle shutdown gracefully
// 	go func() {
// 		sigCh := make(chan os.Signal, 1)
// 		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
// 		<-sigCh
// 		cancel()
// 	}()

// 	// Run application
// 	if err := application.Run(ctx); err != nil {
// 		log.Fatalf("Application error: %v", err)
// 	}
// }

// // package main

// // import (
// //     "log"
// //     "internal/storage"  // Replace with your actual package path
// // )

// // func main() {
// //     cfg := storage.StorageConfig{
// //         Endpoint:        "localhost:9000",
// //         AccessKeyID:     "MINIOADMIN",
// //         SecretAccessKey: "MINIOPASSWORD",
// //         UseSSL:         false,
// //         BucketName:     "models", // Or whatever bucket name you want
// //     }

// //     client, err := storage.NewMinioClient(cfg)
// //     if err != nil {
// //         log.Fatalf("Failed to initialize storage: %v", err)
// //     }

// //     // Now you can use the client to interact with MinIO
// //     // For example:
// //     // Upload a file:
// //     // client.PutObject(context.Background(), cfg.BucketName, "model.pkl", reader, size, minio.PutObjectOptions{})
// // }
