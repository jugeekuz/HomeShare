package helpers

import (
	"os"
	"log"
	"path/filepath"
	"crypto/rand"
	"encoding/base64"
)

func GetOrCreateJWTSecret(dirPath string, fileName string) string {
	filePath := filepath.Join(dirPath, fileName)

	if err := os.MkdirAll(dirPath, 0700); err != nil {
		log.Fatalf("failed to create secrets directory: %v", err)
	}

	if _, err := os.Stat(filePath); err == nil {
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("failed to read secret file: %v", err)
		}
		return string(data)
	} else if !os.IsNotExist(err) {
		log.Fatalf("failed to stat secret file: %v", err)
	}

	secretBytes := make([]byte, 32)
	if _, err := rand.Read(secretBytes); err != nil {
		log.Fatalf("failed to generate secret: %v", err)
	}
	secret := base64.StdEncoding.EncodeToString(secretBytes)

	if err := os.WriteFile(filePath, []byte(secret), 0600); err != nil {
		log.Fatalf("failed to write secret to file: %v", err)
	}

	return secret
}