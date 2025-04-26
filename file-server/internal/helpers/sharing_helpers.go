package helpers

import (
	"fmt"
	"os"
	"time"
	"strings"
	"path/filepath"

)

func GenerateFolderName(expiryDuration time.Duration, folderId string) string {
	timestamp := time.Now().UTC().Truncate(time.Second).Format("20060102150405")

	return fmt.Sprintf("%s_%s_%s", timestamp, expiryDuration, folderId)
}

func CleanupExpiredFolders(root string) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		folderName := entry.Name()
		parts := strings.Split(folderName, "_")
		if len(parts) < 3 {
			fmt.Printf("Skipping folder '%s' (invalid name format)\n", folderName)
			continue
		}

		creationStr := parts[0]
		durationStr := parts[1]

		creationTime, err := time.Parse("20060102150405", creationStr)
		if err != nil {
			fmt.Printf("Error parsing timestamp for folder '%s': %v\n", folderName, err)
			continue
		}

		expiryDuration, err := time.ParseDuration(durationStr)
		if err != nil {
			fmt.Printf("Error parsing duration for folder '%s': %v\n", folderName, err)
			continue
		}

		expiryTime := creationTime.Add(expiryDuration)
		if now.After(expiryTime) {
			folderPath := filepath.Join(root, folderName)
			fmt.Printf("Removing expired folder: %s\n", folderPath)
			if err := os.RemoveAll(folderPath); err != nil {
				fmt.Printf("Error removing folder '%s': %v\n", folderPath, err)
			}
		} else {
			fmt.Printf("Folder '%s' is still active (expires at %s)\n", folderName, expiryTime)
		}
	}

	return nil
}