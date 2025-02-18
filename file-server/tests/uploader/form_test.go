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

func TestParseForm_InvalidFields(t *testing.T) {
	tests := []struct {
		name		string
		fields		map[string]string
		expectedErr	string
	} {
		{
			name: "non numeric chunk index",
			fields: map[string]string{
				"fileId":        uuid.New().String(),
				"fileName":      "testfile",
				"fileExtension": ".txt",
				"md5Hash":       "d41d8cd98f00b204e9800998ecf8427e",
				"chunkIndex":    "not_a_number",
				"totalChunks":   "1",
			},
			expectedErr: "invalid chunk number"
		},
		{
			name: "negative chunk index",
			fields: map[string]string{
				"fileId":        uuid.New().String(),
				"fileName":      "testfile",
				"fileExtension": ".txt",
				"md5Hash":       "d41d8cd98f00b204e9800998ecf8427e",
				"chunkIndex":    "-1",
				"totalChunks":   "1",
			},
			expectedErr: "invalid chunk number"
		},
		{
			name: "chunk index out of range",
			fields: map[string]string{
				"fileId":        uuid.New().String(),
				"fileName":      "testfile",
				"fileExtension": ".txt",
				"md5Hash":       "d41d8cd98f00b204e9800998ecf8427e",
				"chunkIndex":    "1",
				"totalChunks":   "1",
			},
			expectedErr: "invalid chunk number"
		},
		{
			name: "invalid file name",
			fields: map[string]string{
				"fileId":        uuid.New().String(),
				"fileName":      "test/file",
				"fileExtension": ".txt",
				"md5Hash":       "d41d8cd98f00b204e9800998ecf8427e",
				"chunkIndex":    "0",
				"totalChunks":   "1",
			},
			expectedErr: "invalid file name format"
		},
		{
			name: "invalid file id",
			fields: map[string]string{
				"fileId":        "not-a-valid-uuid",
				"fileName":      "test/file",
				"fileExtension": ".txt",
				"md5Hash":       "d41d8cd98f00b204e9800998ecf8427e",
				"chunkIndex":    "0",
				"totalChunks":   "1",
			},
			expectedErr: "invalid file name format"
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := createMultipartRequest(t, tc.fields, "chunk", []byte("content"))
			rr := httptest.NewRecorder()
			_, _, err := uploader.ParseForm(rr, req)
			if err == nil || !strings.Contains(err.Error(), tc.expectedErr) {
				t.Fatalf("expected error containing %q, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestParseForm_InvalidFileExtensions(t *testing.T) {
	invalidFileExtensions := []string{
        // Executable Files
        ".exe",
        ".dll",
        ".bat",
        ".cmd",
        ".com",
        ".msi",
        ".scr",
        ".cpl",
        ".msc",
        // Scripting & Code Files
        ".ps1",
        ".vbs",
        ".wsf",
        ".sh",
        ".php",
        ".php3",
        ".php4",
        ".php5",
        ".phtml",
        ".asp",
        ".aspx",
        ".jsp",
        ".cgi",
        ".pl",
        ".py",
        ".rb",
        ".jar",
        // Other Potentially Unsafe Files
        ".reg",
        ".vbe",
        ".jse",
        ".hta",
        ".lnk",
    }
	tests := make([]struct {
		name        string
		fields      map[string]string
		expectedErr string
	}, len(invalidFileExtensions))

	for i, ext := range invalidFileExtensions {
		tests[i] = struct {
			name        string
			fields      map[string]string
			expectedErr string
		}{
			name: fmt.Sprintf("Testing for extension %s", ext),
			fields: map[string]string{
				"fileId":        uuid.New().String(),
				"fileName":      "testfile",
				"fileExtension": ext,
				"md5Hash":       "d41d8cd98f00b204e9800998ecf8427e",
				"chunkIndex":    "0",
				"totalChunks":   "1",
			},
			expectedErr: fmt.Sprintf("invalid file extension: %s", ext),
		}
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T){
			req := createMultipartRequest(t, tc.fields, "chunk", []byte("content"))
			rr := httptest.NewRecorder()
			_, _, err := uploader.ParseForm(rr, req)
			if err == nil || !string.Contains(err.Error(), tc.expectedErr) {
				t.Fatalf("expected error for invalid file extension, got %v", err)
			}
		})
	}
}
