package helpers

import (
	"os"
	"fmt"
	"io"
	"path/filepath"
	"archive/zip"
	
	"file-server/internal/job"
)


func CreateZip(folderPath string, zipFileName string, files []string, jm *job.JobManager) error {
	jobId := zipFileName

	if !jm.AcquireJob(jobId) {
		return fmt.Errorf("[FILE-SERVER] Error while trying to acquire job Id for zip file : %s", zipFileName)
	}
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
		info, err := os.Stat(filePath)
		if err != nil {
			return err
		}
		if info.IsDir() {
			continue
		}
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
	
	fileName := filepath.Base(filePath)
	writer, err := zipWriter.Create(fileName)
	if err != nil {
		return err
	}

	if _, err = io.Copy(writer, file); err != nil {
		return err
	}

	return nil
}
