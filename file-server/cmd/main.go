// main.go
package main

import (
	// "bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	const MAX_MBYTES = 1
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }

	r.Body = http.MaxBytesReader(w, r.Body, MAX_MBYTES<<20+1024)
	err := r.ParseMultipartForm(MAX_MBYTES << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	
	fileID := r.FormValue("fileID")
	chunkNumberStr := r.FormValue("chunkNumber")
	totalChunksStr := r.FormValue("totalChunks")
	if fileID == "" || chunkNumberStr == "" || totalChunksStr == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}
	if !strings.Contains(fileID, ".") {
		http.Error(w, "FileID doesn't contain an extension", http.StatusBadRequest)
		return
	}
	parts := strings.Split(fileID, ".")
	fileName := parts[0]
	// fileExtension := parts[1]

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
	log.Print(totalChunks)

	file, _, err := r.FormFile("chunk")
	if err != nil {
		http.Error(w, "Error retrieving chunk", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadDir := filepath.Join("uploads", fileName)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, "Error creating directory", http.StatusInternalServerError)
		return
	}

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

	files, err := os.ReadDir(uploadDir)
	if err == nil && len(files) == totalChunks {
		finalFilePath := filepath.Join("uploads", fileID)
		finalFile, err := os.Create(finalFilePath)
		if err != nil {
			http.Error(w, "Error creating final file", http.StatusInternalServerError)
			return
		}
		defer finalFile.Close()

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

		err = os.RemoveAll(uploadDir)
		if err != nil {
			log.Fatalf("Error deleting directory: %v", err)
		}

	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Chunk uploaded successfully"))
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
