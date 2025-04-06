package uploader

import (
	// "bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	
	"file-server/config"
	"file-server/internal/job"
	"file-server/internal/auth"
	
)

type ChunkMeta struct {
	FileId        string
	FileName      string
	FileExtension string
	MD5Hash       string
	ChunkIndex    int
	TotalChunks   int
}

type Chunk struct {
	File multipart.File
}

func getUniqueFileName(path string) string {
	counter := 1
	dir := filepath.Dir(path)
	ext := filepath.Ext(path)
	baseName := strings.TrimSuffix(filepath.Base(path), ext)

	for {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return path
		}

		newName := fmt.Sprintf("%s (%d)%s", baseName, counter, ext)
		path = filepath.Join(dir, newName)
		counter++
	}
}

func ParseFormFileId(w http.ResponseWriter, r *http.Request) (string, error) {
	const MAX_MBYTES = 5

	r.Body = http.MaxBytesReader(w, r.Body, (MAX_MBYTES<<20)+1024)

	if err := r.ParseMultipartForm(MAX_MBYTES << 20); err != nil {
		return "", fmt.Errorf("unable to parse form: %w", err)
	}

	if r.FormValue("fileId") == "" {
		return "", fmt.Errorf("fileId is required")
	}
	return r.FormValue("fileId"), nil
}

func ParseForm(w http.ResponseWriter, r *http.Request) (ChunkMeta, Chunk, error) {
	const MAX_MBYTES = 5

	r.Body = http.MaxBytesReader(w, r.Body, (MAX_MBYTES<<20)+1024)

	if err := r.ParseMultipartForm(MAX_MBYTES << 20); err != nil {
		return ChunkMeta{}, Chunk{}, fmt.Errorf("unable to parse form: %w", err)
	}

	if r.FormValue("fileId") == "" {
		return ChunkMeta{}, Chunk{}, fmt.Errorf("fileId is required")
	}
	if r.FormValue("fileName") == "" {
		return ChunkMeta{}, Chunk{}, fmt.Errorf("fileName is required")
	}
	if r.FormValue("fileExtension") == "" {
		return ChunkMeta{}, Chunk{}, fmt.Errorf("fileExtension is required")
	}
	if r.FormValue("md5Hash") == "" {
		return ChunkMeta{}, Chunk{}, fmt.Errorf("md5Hash is required")
	}

	if r.FormValue("chunkIndex") == "" {
		return ChunkMeta{}, Chunk{}, fmt.Errorf("chunkIndex is required")
	}

	if r.FormValue("totalChunks") == "" {
		return ChunkMeta{}, Chunk{}, fmt.Errorf("totalChunks is required")
	}

	files, ok := r.MultipartForm.File["chunk"]
	if !ok || len(files) == 0 {
		return ChunkMeta{}, Chunk{}, fmt.Errorf("chunk file is required")
	}
	if files[0].Size == 0 {
		return ChunkMeta{}, Chunk{}, fmt.Errorf("chunk file is empty")
	}

	chunkIndex, err := strconv.Atoi(r.FormValue("chunkIndex"))
	if err != nil {
		return ChunkMeta{}, Chunk{}, fmt.Errorf("invalid chunk number: %w", err)
	}

	totalChunks, err := strconv.Atoi(r.FormValue("totalChunks"))
	if err != nil {
		return ChunkMeta{}, Chunk{}, fmt.Errorf("invalid number of chunks: %w", err)
	}

	if chunkIndex > totalChunks-1 || chunkIndex < 0 {
		return ChunkMeta{}, Chunk{}, fmt.Errorf("invalid chunk index: %d", chunkIndex)
	}

	file, _, err := r.FormFile("chunk")
	if err != nil {
		return ChunkMeta{}, Chunk{}, fmt.Errorf("error while reading chunk: %w", err)
	}

	meta := ChunkMeta{
		FileId:        r.FormValue("fileId"),
		FileName:      r.FormValue("fileName"),
		FileExtension: r.FormValue("fileExtension"),
		MD5Hash:       r.FormValue("md5Hash"),
		ChunkIndex:    chunkIndex,
		TotalChunks:   totalChunks,
	}

	fileNameRegex := regexp.MustCompile(`^[a-zA-Z0-9._ -\(\)]+$`)
	if !fileNameRegex.MatchString(meta.FileName) {
		return ChunkMeta{}, Chunk{}, fmt.Errorf("invalid file name format: %s", meta.FileName)
	}

	fileExtensionRegex := regexp.MustCompile(`^\.(jpe?g|png|gif|bmp|tiff?|webp|mp4|mkv|mov|avi|flv|wmv|txt)$`)
	if !fileExtensionRegex.MatchString(strings.ToLower(meta.FileExtension)) {
		return ChunkMeta{}, Chunk{}, fmt.Errorf("invalid file extension: %s", meta.FileExtension)
	}

	if _, err := uuid.Parse(meta.FileId); err != nil {
		return ChunkMeta{}, Chunk{}, fmt.Errorf("invalid UUID format: %w", err)
	}

	chunk := Chunk{
		File: file,
	}

	return meta, chunk, nil
}

func ChunkAssemble(meta ChunkMeta, jm *job.JobManager, absolutePath string) {
	cfg := config.LoadConfig()

	defer jm.ReleaseJob(meta.FileId)

	chunksDir := filepath.Join(absolutePath, cfg.ChunksDir, meta.FileId)
	if _, err := os.Stat(chunksDir); os.IsNotExist(err) {
		log.Printf("Chunk directory %s does not exist for file ID: %s", chunksDir, meta.FileId)
		return
	}

	defer func() {
		if err := os.RemoveAll(chunksDir); err != nil {
			log.Printf("Error deleting chunk directory %s: %v", chunksDir, err)
		}
	}()

	finalFilePath := filepath.Join(absolutePath, meta.FileName+meta.FileExtension)

	finalFilePath = getUniqueFileName(finalFilePath) // If file exists then save as `file (1)`

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
		finalFile.Close() // File has to be closed in order to be removed

		log.Printf("MD5 mismatch for %s. Computed: %s, Expected: %s", meta.FileId, computedHash, expectedHash)
		if err := os.Remove(finalFilePath); err != nil {
			log.Printf("Error while deleting final file path %s: %q", finalFilePath, err)
		}

		return
	}

	log.Printf("Successfully assembled file %s", finalFilePath)
}

func UploadHandler(w http.ResponseWriter, r *http.Request, jm *job.JobManager, folderPath string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Check if user has write access to the folder
	claimsRaw := r.Context().Value(auth.ClaimsContextKey)
	claims, ok := claimsRaw.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	canAccess, err := auth.HasAccess(claims, filepath.Base(folderPath), "w")
	if err != nil || !canAccess {
		http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
		return
	}
	
	cfg := config.LoadConfig()

	meta, chunk, err := ParseForm(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer chunk.File.Close()

	chunksDir := filepath.Join(folderPath, cfg.ChunksDir, meta.FileId)
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

	files, err := os.ReadDir(chunksDir)
	if err != nil {
		http.Error(w, "Error reading chunk directory", http.StatusInternalServerError)
		return
	}

	if len(files) == meta.TotalChunks {
		if (jm.AcquireJob(meta.FileId)) {
			go ChunkAssemble(meta, jm, folderPath)
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Chunk uploaded successfully"))
}