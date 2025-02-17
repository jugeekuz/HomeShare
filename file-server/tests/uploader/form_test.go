package uploader_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"

	"file-server/internal/uploader"
)

// createMultipartRequest is a helper to build a multipart HTTP request.
func createMultipartRequest(t *testing.T, fields map[string]string, fileField string, fileContent []byte) *http.Request {
	t.Helper()

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Add form fields.
	for k, v := range fields {
		if err := w.WriteField(k, v); err != nil {
			t.Fatalf("error writing field %s: %v", k, err)
		}
	}

	// Add the file if specified.
	if fileField != "" {
		fw, err := w.CreateFormFile(fileField, "dummy")
		if err != nil {
			t.Fatalf("error creating form file: %v", err)
		}
		if _, err := fw.Write(fileContent); err != nil {
			t.Fatalf("error writing file content: %v", err)
		}
	}

	if err := w.Close(); err != nil {
		t.Fatalf("error closing writer: %v", err)
	}

	req := httptest.NewRequest("POST", "/upload", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req
}

//
// ParseForm Tests
//

func TestParseForm_MissingFields(t *testing.T) {
	// Table-driven tests for missing required fields.
	tests := []struct {
		name        string
		fields      map[string]string
		expectedErr string
	}{
		{
			name: "missing fileId",
			fields: map[string]string{
				"fileName":      "testfile",
				"fileExtension": ".txt",
				"md5Hash":       "d41d8cd98f00b204e9800998ecf8427e",
				"chunkIndex":    "0",
				"totalChunks":   "1",
			},
			expectedErr: "fileId is required",
		},
		{
			name: "missing fileName",
			fields: map[string]string{
				"fileId":        uuid.New().String(),
				"fileExtension": ".txt",
				"md5Hash":       "d41d8cd98f00b204e9800998ecf8427e",
				"chunkIndex":    "0",
				"totalChunks":   "1",
			},
			expectedErr: "fileName is required",
		},
		{
			name: "missing fileExtension",
			fields: map[string]string{
				"fileName":      "testfile",
				"fileId":        uuid.New().String(),
				"md5Hash":       "d41d8cd98f00b204e9800998ecf8427e",
				"chunkIndex":    "0",
				"totalChunks":   "1",
			},
			expectedErr: "fileExtension is required",
		},
		{
			name: "missing md5Hash",
			fields: map[string]string{
				"fileName":      "testfile",
				"fileId":        uuid.New().String(),
				"fileExtension": ".txt",
				"chunkIndex":    "0",
				"totalChunks":   "1",
			},
			expectedErr: "md5Hash is required",
		},
		{
			name: "missing chunkIndex",
			fields: map[string]string{
				"fileName":      "testfile",
				"fileId":        uuid.New().String(),
				"fileExtension": ".txt",
				"md5Hash":       "d41d8cd98f00b204e9800998ecf8427e",
				"totalChunks":   "1",
			},
			expectedErr: "chunkIndex is required",
		},
		{
			name: "missing chunkIndex",
			fields: map[string]string{
				"fileName":      "testfile",
				"fileId":        uuid.New().String(),
				"fileExtension": ".txt",
				"md5Hash":       "d41d8cd98f00b204e9800998ecf8427e",
				"chunkIndex":    "0",
			},
			expectedErr: "totalChunks is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := createMultipartRequest(t, tc.fields, "chunk", []byte("dummy content"))
			rr := httptest.NewRecorder()
			_, _, err := uploader.ParseForm(rr, req)
			if err == nil || !strings.Contains(err.Error(), tc.expectedErr) {
				t.Fatalf("expected error containing %q, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestParseForm_InvalidChunkIndex(t *testing.T) {
	fields := map[string]string{
		"fileId":        uuid.New().String(),
		"fileName":      "testfile",
		"fileExtension": ".txt",
		"md5Hash":       "d41d8cd98f00b204e9800998ecf8427e",
		"chunkIndex":    "not_a_number",
		"totalChunks":   "1",
	}
	req := createMultipartRequest(t, fields, "chunk", []byte("content"))
	rr := httptest.NewRecorder()
	_, _, err := uploader.ParseForm(rr, req)
	if err == nil || !strings.Contains(err.Error(), "invalid chunk number") {
		t.Fatalf("expected error for invalid chunk number, got %v", err)
	}
}

func TestParseForm_InvalidChunkIndexOutOfRange(t *testing.T) {
	fields := map[string]string{
		"fileId":        uuid.New().String(),
		"fileName":      "testfile",
		"fileExtension": ".txt",
		"md5Hash":       "d41d8cd98f00b204e9800998ecf8427e",
		"chunkIndex":    "2",
		"totalChunks":   "2", 
	}
	req := createMultipartRequest(t, fields, "chunk", []byte("content"))
	rr := httptest.NewRecorder()
	_, _, err := uploader.ParseForm(rr, req)
	if err == nil || !strings.Contains(err.Error(), "invalid chunk index") {
		t.Fatalf("expected error for out-of-range chunk index, got %v", err)
	}
}

func TestParseForm_InvalidFileName(t *testing.T) {
	fields := map[string]string{
		"fileId":        uuid.New().String(),
		"fileName":      "test/file",
		"fileExtension": ".txt",
		"md5Hash":       "d41d8cd98f00b204e9800998ecf8427e",
		"chunkIndex":    "0",
		"totalChunks":   "1",
	}
	req := createMultipartRequest(t, fields, "chunk", []byte("content"))
	rr := httptest.NewRecorder()
	_, _, err := uploader.ParseForm(rr, req)
	if err == nil || !strings.Contains(err.Error(), "invalid file name format") {
		t.Fatalf("expected error for invalid file name, got %v", err)
	}
}

func TestParseForm_InvalidFileExtension(t *testing.T) {
	// TODO: test for more file extensions
	fields := map[string]string{
		"fileId":        uuid.New().String(),
		"fileName":      "testfile",
		"fileExtension": ".exe",
		"md5Hash":       "d41d8cd98f00b204e9800998ecf8427e",
		"chunkIndex":    "0",
		"totalChunks":   "1",
	}
	req := createMultipartRequest(t, fields, "chunk", []byte("content"))
	rr := httptest.NewRecorder()
	_, _, err := uploader.ParseForm(rr, req)
	if err == nil || !strings.Contains(err.Error(), "invalid file extension") {
		t.Fatalf("expected error for invalid file extension, got %v", err)
	}
}

func TestParseForm_InvalidUUID(t *testing.T) {
	fields := map[string]string{
		"fileId":        "not-a-valid-uuid",
		"fileName":      "testfile",
		"fileExtension": ".txt",
		"md5Hash":       "d41d8cd98f00b204e9800998ecf8427e",
		"chunkIndex":    "0",
		"totalChunks":   "1",
	}
	req := createMultipartRequest(t, fields, "chunk", []byte("content"))
	rr := httptest.NewRecorder()
	_, _, err := uploader.ParseForm(rr, req)
	if err == nil || !strings.Contains(err.Error(), "invalid UUID format") {
		t.Fatalf("expected error for invalid UUID, got %v", err)
	}
}

func TestParseForm_Valid(t *testing.T) {
	fields := map[string]string{
		"fileId":        uuid.New().String(),
		"fileName":      "testfile",
		"fileExtension": ".txt",
		"md5Hash":     "098f6bcd4621d373cade4e832627b4f6",
		"chunkIndex":  "0",
		"totalChunks": "1",
	}
	req := createMultipartRequest(t, fields, "chunk", []byte("test"))
	rr := httptest.NewRecorder()
	meta, chunk, err := uploader.ParseForm(rr, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if meta.FileName != "testfile" || meta.FileExtension != ".txt" || meta.ChunkIndex != 0 || meta.TotalChunks != 1 {
		t.Errorf("metadata not parsed correctly: %+v", meta)
	}
	if chunk.File == nil {
		t.Error("expected valid file in chunk")
	}
}
