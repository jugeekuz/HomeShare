// main.go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse a multipart form (adjust maxMemory as needed)
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve form values
	fileID := r.FormValue("fileID")
	chunkNumberStr := r.FormValue("chunkNumber")
	totalChunksStr := r.FormValue("totalChunks")
	if fileID == "" || chunkNumberStr == "" || totalChunksStr == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	chunkNumber, err := strconv.Atoi(chunkNumberStr)
	if err != nil {
		http.Error(w, "Invalid chunk number", http.StatusBadRequest)
		return
	}
	totalChunks, err := strconv.Atoi(totalChunksStr)
	if err != nil {
		http.Error(w, "Invalid total chunks", http.StatusBadRequest)
		return
	}

	// Retrieve the file chunk from the form
	file, _, err := r.FormFile("chunk")
	if err != nil {
		http.Error(w, "Error retrieving chunk", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a temporary directory for the file's chunks (if it doesn't exist)
	uploadDir := filepath.Join("uploads", fileID)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, "Error creating directory", http.StatusInternalServerError)
		return
	}

	// Save the chunk to disk
	chunkFilePath := filepath.Join(uploadDir, fmt.Sprintf("chunk_%d", chunkNumber))
	out, err := os.Create(chunkFilePath)
	if err != nil {
		http.Error(w, "Unable to create chunk file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Error saving chunk", http.StatusInternalServerError)
		return
	}

	// Optionally, check if all chunks have been uploaded
	files, err := os.ReadDir(uploadDir)
	if err == nil && len(files) == totalChunks {
		finalFilePath := filepath.Join("uploads", fmt.Sprintf("%s_final", fileID))
		finalFile, err := os.Create(finalFilePath)
		if err != nil {
			http.Error(w, "Error creating final file", http.StatusInternalServerError)
			return
		}
		defer finalFile.Close()

		// Reassemble chunks in order
		for i := 1; i <= totalChunks; i++ {
			chunkPath := filepath.Join(uploadDir, fmt.Sprintf("chunk_%d", i))
			chunkFile, err := os.Open(chunkPath)
			if err != nil {
				http.Error(w, "Error opening chunk", http.StatusInternalServerError)
				return
			}
			_, err = io.Copy(finalFile, chunkFile)
			chunkFile.Close()
			if err != nil {
				http.Error(w, "Error assembling chunks", http.StatusInternalServerError)
				return
			}
		}
		// Optionally: Remove chunk files after assembly.
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Chunk uploaded successfully"))
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
