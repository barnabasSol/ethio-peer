package common

import (
	"log"
	"os"
)

type MinioConfig struct {
	Key         string
	Secret      string
	URL         string
	ImageBucket string
}

func GetMinioConfig() MinioConfig {
	config := MinioConfig{
		Key:         os.Getenv("MINIO_KEY"),
		Secret:      os.Getenv("MINIO_SECRET"),
		URL:         os.Getenv("MINIO_URL"),
		ImageBucket: os.Getenv("MINIO_IMAGE_BUCKET"),
	}

	if config.Key == "" || config.Secret == "" || config.URL == "" || config.ImageBucket == "" {
		log.Fatal("One or more MINIO environment variables are not set")
	}

	return config
}
