package helpers

import (
	"os"
	"io"
	"path/filepath"
	"archive/zip"
	
	"file-server/internal/job"
)


func CreateZip(folderPath string, zipFileName string, files []string, jm *job.JobManager) error {
	jobId := zipFileName

	jm.AcquireJob(jobId)
	defer jm.ReleaseJob(jobId)

	zipFilePath := filepath.Join(folderPath, zipFileName)
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range files {
		filePath := filepath.Join(folderPath, file)
		if err := AddFileToZip(zipWriter, filePath); err != nil {
			return err
		}
	}

	return nil
}

// AddFileToZip adds an individual file to the zip archive and logs each step.
func AddFileToZip(zipWriter *zip.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer, err := zipWriter.Create(filePath)
	if err != nil {
		return err
	}

	if _, err = io.Copy(writer, file); err != nil {
		return err
	}

	return nil
}
