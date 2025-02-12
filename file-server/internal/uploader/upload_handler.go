package uploader

import (
	// "bufio"
	"fmt"
	"log"
	"os"
	"io"
	"regexp"
	"strconv"
	"strings"
	"crypto/md5"
	"encoding/hex"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"github.com/google/uuid"
)

type ChunkMeta struct {
	FileId			string
    FileName    	string
	FileExtension 	string
    MD5Hash     	string
    ChunkIndex  	int
    TotalChunks 	int
}

type Chunk struct {
	File multipart.File
}

var (
	UploadDir = "uploads"
)
var (
	ChunkDir = "chunks"
)

func parseForm(w http.ResponseWriter, r *http.Request) (ChunkMeta, Chunk, error) {
    const MAX_MBYTES = 1

    r.Body = http.MaxBytesReader(w, r.Body, (MAX_MBYTES<<20)+1024)

    if err := r.ParseMultipartForm(MAX_MBYTES << 20); err != nil {
		return ChunkMeta{}, Chunk{}, fmt.Errorf("unable to parse form: %w", err)
	}

    chunkIndex, err := strconv.Atoi(r.FormValue("chunkIndex"))
	if err != nil {
		return ChunkMeta{}, Chunk{}, fmt.Errorf("invalid chunk number: %w", err)
	}

	totalChunks, err := strconv.Atoi(r.FormValue("totalChunks"))
	if err != nil {
		return ChunkMeta{}, Chunk{}, fmt.Errorf("invalid number of chunks: %w", err)
	}

	file, _, err := r.FormFile("chunk")
	if err != nil {
        return ChunkMeta{}, Chunk{}, fmt.Errorf("error while reading chunk: %w", err)
	}

    meta := ChunkMeta{
        FileId:    		r.FormValue("FileId"),
        FileName:    	r.FormValue("FileName"),
		FileExtension: 	r.FormValue("FileExtension"),
        MD5Hash:     	r.FormValue("MD5Hash"),
        ChunkIndex:  	chunkIndex,
		TotalChunks: 	totalChunks,
    }

	fileNameRegex := regexp.MustCompile(`^[a-zA-Z0-9._ -]+$`)
	if !fileNameRegex.MatchString(meta.FileName) {
		return ChunkMeta{}, Chunk{}, fmt.Errorf("invalid file name format for file ID: %s", meta.FileName)
	}

	fileExtensionRegex := regexp.MustCompile(`^\.(jpe?g|png|gif|bmp|tiff?|webp|mp4|mkv|mov|avi|flv|wmv)$`)
	if !fileExtensionRegex.MatchString(meta.FileExtension) {
		return ChunkMeta{}, Chunk{}, fmt.Errorf("invalid file extension format for file ID: %s", meta.FileExtension)
	}

	
    if _, err := uuid.Parse(meta.FileId); err != nil {
        return ChunkMeta{}, Chunk{}, fmt.Errorf("invalid UUID format for file ID: %w", err)
    }

	chunk := Chunk{
		File: file,
	}

    return meta, chunk, nil
}

func chunkAssemble(meta ChunkMeta) {
	chunksDir := filepath.Join(ChunkDir, meta.FileId)
    if _, err := os.Stat(chunksDir); os.IsNotExist(err) {
        log.Printf("Chunk directory %s does not exist for file ID: %s", chunksDir, meta.FileId)
        return
    }

	defer func () {
		if err := os.RemoveAll(chunksDir); err != nil {
			log.Printf("Error deleting chunk directory %s: %v", chunksDir, err)
		}
	}()

    files, err := os.ReadDir(chunksDir)
    if err != nil {
        log.Printf("Error reading chunk directory %s: %v", chunksDir, err)
        return
    }
    if len(files) != meta.TotalChunks {
        log.Printf("Incomplete chunks for file ID %s: expected %d, found %d", meta.FileId, meta.TotalChunks, len(files))
        return
    }

	finalFilePath := filepath.Join(UploadDir, meta.FileName + meta.FileExtension)
	finalFile, err := os.Create(finalFilePath)
	if err != nil {
		log.Printf("Error creating final file %s: %v", finalFilePath, err)
        return
	}
	defer finalFile.Close()

	hasher := md5.New()

	for i := 0; i < meta.TotalChunks; i++ {
		chunkPath := filepath.Join(chunksDir, fmt.Sprintf("chunk_%d", i))
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			log.Printf("Error while opening chunk %s : %q", chunkPath, err)
			return
		}

		multiWriter := io.MultiWriter(finalFile, hasher)
        if _, err := io.Copy(multiWriter, chunkFile); err != nil {
            chunkFile.Close()
            log.Printf("Error copying chunk %s: %v", chunkPath, err)
            return
        }
		chunkFile.Close()
	}

	computedHash := hex.EncodeToString(hasher.Sum(nil))
    expectedHash := strings.ToLower(strings.TrimSpace(meta.MD5Hash))

    if strings.ToLower(computedHash) != expectedHash {
        log.Printf("MD5 mismatch for %s. Computed: %s, Expected: %s", meta.FileId, computedHash, expectedHash)
        os.Remove(finalFilePath)
        return
    }
	
    log.Printf("Successfully assembled file %s", finalFilePath)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	meta, chunk, err := parseForm(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer chunk.File.Close()
	
	chunksDir := filepath.Join(ChunkDir, meta.FileId)
	if err := os.MkdirAll(chunksDir, os.ModePerm); err != nil {
		http.Error(w, "Error creating directory", http.StatusInternalServerError)
		return
	}

	chunkFilePath := filepath.Join(chunksDir, fmt.Sprintf("chunk_%d", meta.ChunkIndex))
	out, err := os.Create(chunkFilePath)
	if err != nil {
		http.Error(w, "Unable to create chunk file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, chunk.File)
	if err != nil {
		http.Error(w, "Error saving chunk", http.StatusInternalServerError)
		return
	}

	go chunkAssemble(meta)
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Chunk uploaded successfully"))
}
