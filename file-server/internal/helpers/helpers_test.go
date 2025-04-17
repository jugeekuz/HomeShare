package helpers

import (
	"archive/zip"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"file-server/internal/job"

	"github.com/google/uuid"
)

var TestFolder string = GenerateFolderName(0 * time.Second, uuid.New().String())

func TestMain(m *testing.M) {
	if err := os.MkdirAll(TestFolder, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create testing directory %s: %v\n", TestFolder, err)
		os.Exit(1)
	}

	exitCode := m.Run()

	if err := os.RemoveAll(TestFolder); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to remove testing directory %s: %v\n", TestFolder, err)
	}

	os.Exit(exitCode)
}

func TestZipFileCreation(t *testing.T) {
	zipFileName := "zip_file_test.zip"
	bufferSize := 5*1024*1024 // 5MB
	buf := make([]byte, bufferSize)

	// Create File 1
	file1 := "someFile1.txt"
	filePath1 := filepath.Join(TestFolder, file1)
	out, err := os.Create(filePath1)
	if err != nil {
		t.Fatalf("Received unexpected error when creating file: %v", err)
	}
	if _, err := rand.Read(buf); err != nil {
		t.Fatalf("Received unexpected error when creating random bytes: %v", err)
	}
	if _, err := out.Write(buf); err != nil {
		t.Fatalf("Received unexpected error when writing into file: %v", err)
	}
	out.Close()

	// Create File 2
	file2 := "someFile2.txt"
	filePath2 := filepath.Join(TestFolder, file2)
	out, err = os.Create(filePath2)
	if err != nil {
		t.Fatalf("Received unexpected error when creating file: %v", err)
	}
	if _, err := rand.Read(buf); err != nil {
		t.Fatalf("Received unexpected error when creating random bytes: %v", err)
	}
	if _, err := out.Write(buf); err != nil {
		t.Fatalf("Received unexpected error when writing into file: %v", err)
	}
	out.Close()

	// Create File 3
	file3 := "someFile3.txt"
	filePath3 := filepath.Join(TestFolder, file3)
	out, err = os.Create(filePath3)
	if err != nil {
		t.Fatalf("Received unexpected error when creating file: %v", err)
	}
	if _, err := rand.Read(buf); err != nil {
		t.Fatalf("Received unexpected error when creating random bytes: %v", err)
	}
	if _, err := out.Write(buf); err != nil {
		t.Fatalf("Received unexpected error when writing into file: %v", err)
	}
	out.Close()

	files := make([]string, 2)
	files[0] = file1
	files[1] = file2

	jm := job.NewJobManager(30 * time.Minute)

	// Test only the two files in the beggining
	t.Run("Test_Zip_Initial_Creation", func(t *testing.T) {
		if err := CreateZip(TestFolder, zipFileName, files, jm); err != nil {
			t.Errorf("Received Unexpected error when creating zip file: %v", err)
		}

		zipFilePath := filepath.Join(TestFolder, zipFileName)
		if _, err := os.Stat(zipFilePath); err != nil {
			if os.IsNotExist(err) {
				t.Errorf("Zip file does not exist")
			} else {
				t.Errorf("Received unexpected error when reading zip file: %v", err)
			}
		}

		reader, err := zip.OpenReader(zipFilePath)
		if err != nil {
			t.Fatalf("Received unexpected error when reading zipFile: %v", err)
		}
		defer reader.Close()

		if len(reader.File) != 2 {
			t.Errorf("Expected zip file length 2, received: %d", len(reader.File))
		}

		for _, file := range reader.File {
			if file.UncompressedSize64 != uint64(bufferSize) {
				t.Errorf("Expected zip file size %d, received: %d", bufferSize, file.UncompressedSize64)
			}
		}
	})

	files = make([]string, 3)
	files[0] = file1
	files[1] = file2
	files[2] = file3

	t.Run("Test_Zip_Adding_Another_File", func(t *testing.T) {
		if err := CreateZip(TestFolder, zipFileName, files, jm); err != nil {
			t.Errorf("Received Unexpected error when creating zip file: %v", err)
		}

		zipFilePath := filepath.Join(TestFolder, zipFileName)
		if _, err := os.Stat(zipFilePath); err != nil {
			if os.IsNotExist(err) {
				t.Errorf("Zip file does not exist")
			} else {
				t.Errorf("Received unexpected error when reading zip file: %v", err)
			}
		}

		reader, err := zip.OpenReader(zipFilePath)
		if err != nil {
			t.Fatalf("Received unexpected error when reading zipFile: %v", err)
		}
		defer reader.Close()

		if len(reader.File) != 3 {
			t.Errorf("Expected zip file length 3, received: %d", len(reader.File))
		}

		for _, file := range reader.File {
			if file.UncompressedSize64 != uint64(bufferSize) {
				t.Errorf("Expected zip file size %d, received: %d", bufferSize, file.UncompressedSize64)
			}
		}
	})

}


func TestZipFileCleanup(t *testing.T) {
	currentDir, err := os.Getwd()
    if err != nil {
        t.Fatalf("Error getting current directory: %v", err)
    }
    
	if err := CleanupExpiredFolders(currentDir); err != nil {
		t.Errorf("Received unexpected error when cleaning expired folders: %v", err)
	}

	if _, err := os.Stat(TestFolder); err == nil {
		t.Errorf("Folder exists, expected for it to be deleted")
	} else if !os.IsNotExist(err) {
		t.Errorf("Received unexpected error when checking folder: %v", err)
	}
}