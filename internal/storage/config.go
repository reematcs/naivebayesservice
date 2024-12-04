package storage

type Config struct {
	Server  ServerConfig
	Storage StorageConfig
}

type ServerConfig struct {
	Port int
	Host string
}

type StorageConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
}

func Load() (*Config, error) {
	// Load configuration from environment or file
	// For now, return default values
	return &Config{
		Server: ServerConfig{
			Port: 8080,
			Host: "localhost",
		},
		Storage: StorageConfig{
			Endpoint:        "localhost:9000",
			AccessKeyID:     "MINIOADMIN",
			SecretAccessKey: "MINIOPASSWORD",
			UseSSL:          false,
			BucketName:      "models",
		},
	}, nil
}
