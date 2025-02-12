package app_test

// import (
// 	"bytes"
// 	"io"
// 	"mime/multipart"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"path/filepath"
// 	"testing"
// 	"file-server/internal/upload_handler"
	
// )

// func TestUploadHandler_OptionsRequest(t *testing.T) {
// 	req := httptest.NewRequest("OPTIONS", "/upload", nil)
// 	rr := httptest.NewRecorder()
// 	uploader.UploadHandler(rr, req)

// 	if rr.Code != http.StatusOK {
// 		t.Errorf("Expected status 200, got %d", rr.Code)
// 	}
// }

// func TestUploadHandler_SuccessfulUpload(t *testing.T) {
// 	// Create temporary directory for test isolation
// 	tempDir := t.TempDir()
// 	oldUploads := uploader.UploadDir
// 	uploader.UploadDir = filepath.Join(tempDir, "uploads")
// 	defer func() { uploader.UploadDir = oldUploads }()

// 	// Create multipart form
// 	body := &bytes.Buffer{}
// 	writer := multipart.NewWriter(body)
// 	writer.WriteField("fileID", "test.txt")
// 	writer.WriteField("chunkNumber", "1")
// 	writer.WriteField("totalChunks", "1")
// 	part, _ := writer.CreateFormFile("chunk", "test.txt")
// 	io.WriteString(part, "test content")
// 	writer.Close()

// 	req := httptest.NewRequest("POST", "/upload", body)
// 	req.Header.Set("Content-Type", writer.FormDataContentType())
// 	rr := httptest.NewRecorder()

// 	uploader.UploadHandler(rr, req)

// 	if rr.Code != http.StatusOK {
// 		t.Errorf("Expected status 200, got %d", rr.Code)
// 	}

// 	// Verify file creation
// 	finalPath := filepath.Join(tempDir, "uploads", "test.txt")
// 	if _, err := os.Stat(finalPath); os.IsNotExist(err) {
// 		t.Errorf("Final file not created")
// 	}
// }