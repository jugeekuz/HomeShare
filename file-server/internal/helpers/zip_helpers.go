package helpers

import (
	"os"
	"io"
	"log"
	"path/filepath"
	"archive/zip"
	
	"file-server/internal/job"
)


func CreateZip(folderPath string, zipFileName string, files []string, jm *job.JobManager) error {
	jobId := zipFileName
	log.Printf("Starting job %s", jobId)
	jm.AcquireJob(jobId)
	defer func() {
		jm.ReleaseJob(jobId)
		log.Printf("Released job %s", jobId)
	}()

	zipFilePath := filepath.Join(folderPath, zipFileName)
	log.Printf("Creating zip file at %s", zipFilePath)
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		log.Printf("Error creating zip file: %v", err)
		return err
	}
	defer func() {
		if cerr := zipFile.Close(); cerr != nil {
			log.Printf("Error closing zip file: %v", cerr)
		} else {
			log.Printf("Zip file closed successfully")
		}
	}()

	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		if cerr := zipWriter.Close(); cerr != nil {
			log.Printf("Error closing zip writer: %v", cerr)
		} else {
			log.Printf("Zip writer closed successfully")
		}
	}()

	for _, file := range files {
		log.Printf("Adding file %s to zip", file)
		filePath := filepath.Join(folderPath, file)
		if err := AddFileToZip(zipWriter, filePath); err != nil {
			log.Printf("Error adding file %s: %v", file, err)
			return err
		}
		log.Printf("Added file %s successfully", file)
	}

	log.Printf("Zip creation completed successfully")
	return nil
}

// AddFileToZip adds an individual file to the zip archive and logs each step.
func AddFileToZip(zipWriter *zip.Writer, filePath string) error {
	log.Printf("Opening file %s", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file %s: %v", filePath, err)
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Printf("Error closing file %s: %v", filePath, cerr)
		} else {
			log.Printf("File %s closed successfully", filePath)
		}
	}()

	log.Printf("Creating zip entry for %s", filePath)
	writer, err := zipWriter.Create(filePath)
	if err != nil {
		log.Printf("Error creating zip entry for %s: %v", filePath, err)
		return err
	}

	log.Printf("Copying contents of %s into zip", filePath)
	if _, err = io.Copy(writer, file); err != nil {
		log.Printf("Error copying file %s: %v", filePath, err)
		return err
	}

	log.Printf("File %s added to zip", filePath)
	return nil
}
