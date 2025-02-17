package uploader_test

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"encoding/hex"
	"strings"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"

	"file-server/internal/uploader"
)


func TestUploadHandler_Valid(t *testing.T) {
	uploadDir := t.TempDir()
	chunkDir := t.TempDir()
	uploader.UploadDir = uploadDir
	uploader.ChunkDir = chunkDir

	fields := map[string]string{
		"fileId":        uuid.New().String(),
		"fileName":      "testfile",
		"fileExtension": ".txt",
		"md5Hash":     fmt.Sprintf("%x", md5.Sum([]byte("chunkdata"))),
		"chunkIndex":  "0",
		"totalChunks": "1",
	}
	req := createMultipartRequest(t, fields, "chunk", []byte("chunkdata"))
	rr := httptest.NewRecorder()

	// Call the handler.
	uploader.UploadHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	chunksPath := filepath.Join(chunkDir, fields["fileId"], "chunk_0")
	if _, err := os.Stat(chunksPath); err != nil {
		t.Errorf("expected chunk file %s to exist, but got error: %v", chunksPath, err)
	}
}

//
// ChunkAssemble Tests
//

func TestChunkAssemble_NoChunkDir(t *testing.T) {
	// If the chunk directory does not exist, no final file should be created.
	chunkDir := t.TempDir()
	uploadDir := t.TempDir()
	uploader.ChunkDir = filepath.Join(chunkDir, "nonexistent")
	uploader.UploadDir = uploadDir

	meta := uploader.ChunkMeta{
		FileId:        uuid.New().String(),
		FileName:      "testfile",
		FileExtension: ".txt",
		// md5 of "dummy"
		MD5Hash:     fmt.Sprintf("%x", md5.Sum([]byte("dummy"))),
		ChunkIndex:  0,
		TotalChunks: 1,
	}
	uploader.ChunkAssemble(meta)
	finalFilePath := filepath.Join(uploadDir, meta.FileName+meta.FileExtension)
	if _, err := os.Stat(finalFilePath); !os.IsNotExist(err) {
		t.Errorf("final file should not exist as chunk directory doesn't exist")
	}
}

func TestChunkAssemble_IncompleteChunks(t *testing.T) {
	// If not all expected chunks are present, no file should be assembled.
	chunkDir := t.TempDir()
	uploadDir := t.TempDir()
	fileId := uuid.New().String()
	chunksDir := filepath.Join(chunkDir, fileId)
	if err := os.MkdirAll(chunksDir, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	// Write only one chunk while expecting two.
	if err := os.WriteFile(filepath.Join(chunksDir, "chunk_0"), []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	uploader.ChunkDir = chunkDir
	uploader.UploadDir = uploadDir

	meta := uploader.ChunkMeta{
		FileId:        fileId,
		FileName:      "testfile",
		FileExtension: ".txt",
		MD5Hash:       fmt.Sprintf("%x", md5.Sum([]byte("data"))),
		ChunkIndex:    0,
		TotalChunks:   2,
	}
	uploader.ChunkAssemble(meta)
	finalFilePath := filepath.Join(uploadDir, meta.FileName+meta.FileExtension)
	if _, err := os.Stat(finalFilePath); !os.IsNotExist(err) {
		t.Errorf("final file should not be assembled because not all chunks are present")
	}
}

func TestChunkAssemble_Success(t *testing.T) {
	chunkDir := t.TempDir()
	uploadDir := t.TempDir()
	fileId := uuid.New().String()
	chunksDir := filepath.Join(chunkDir, fileId)
	if err := os.MkdirAll(chunksDir, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	chunkContents := []string{"hello ", "world"}
	totalData := strings.Join(chunkContents, "")
	hash := md5.Sum([]byte(totalData))

	expectedHash := hex.EncodeToString(hash[:])

	for i, content := range chunkContents {
		chunkPath := filepath.Join(chunksDir, fmt.Sprintf("chunk_%d", i))
		if err := os.WriteFile(chunkPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	uploader.ChunkDir = chunkDir
	uploader.UploadDir = uploadDir

	meta := uploader.ChunkMeta{
		FileId:        fileId,
		FileName:      "assembled",
		FileExtension: ".txt",
		MD5Hash:       expectedHash,
		ChunkIndex:    0,
		TotalChunks:   len(chunkContents),
	}
	uploader.ChunkAssemble(meta)

	finalFilePath := filepath.Join(uploadDir, meta.FileName+meta.FileExtension)
	data, err := os.ReadFile(finalFilePath)
	if err != nil {
		t.Fatalf("expected final file to exist: %v", err)
	}
	if string(data) != totalData {
		t.Errorf("final file content mismatch: got %q, expected %q", string(data), totalData)
	}
}

func TestChunkAssemble_MD5Mismatch(t *testing.T) {
	chunkDir := t.TempDir()
	uploadDir := t.TempDir()
	fileId := uuid.New().String()
	chunksDir := filepath.Join(chunkDir, fileId)
	if err := os.MkdirAll(chunksDir, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	chunkContents := []string{"hello ", "world"}
	wrongHash := "00000000000000000000000000000000"

	for i, content := range chunkContents {
		chunkPath := filepath.Join(chunksDir, fmt.Sprintf("chunk_%d", i))
		if err := os.WriteFile(chunkPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	uploader.ChunkDir = chunkDir
	uploader.UploadDir = uploadDir

	meta := uploader.ChunkMeta{
		FileId:        fileId,
		FileName:      "assembled",
		FileExtension: ".txt",
		MD5Hash:       wrongHash,
		ChunkIndex:    0,
		TotalChunks:   len(chunkContents),
	}
	uploader.ChunkAssemble(meta)

	finalFilePath := filepath.Join(uploadDir, meta.FileName+meta.FileExtension)
	if _, err := os.Stat(finalFilePath); !os.IsNotExist(err) {
		t.Errorf("final file should not exist due to MD5 mismatch")
	}
}
