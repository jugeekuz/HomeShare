package config

import (
	"os"
)

type Config struct {
	UploadDir         string
	SharingDir        string
}

func LoadConfig() *Config {
	return &Config{
		UploadDir: 		getEnv("UPLOAD_DIR", "/uploads"),
		SharingDir: 	getEnv("SHARING_DIR", "/temp"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
