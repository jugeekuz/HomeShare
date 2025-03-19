package uploader

import (
	"bytes"
	"context"
	"crypto/md5"
    "encoding/hex"
	"crypto/rand"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"

	"file-server/config"
	"file-server/internal/job"
	"file-server/internal/auth"
)

type FormFields struct {
	fileId 			string
	fileName 		string
	fileExtension 	string
	md5Hash 		string
	chunkIndex 		string
	totalChunks 	string
	chunkContent 	[]byte
}

// --------------------------------------
// 			 Helper Functions
// --------------------------------------
func createMultipartForm(formFields FormFields) (*http.Request, error) {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	if formFields.fileId != "" {
		if err := writer.WriteField("fileId", formFields.fileId); err != nil {
			return nil, err
		}
	}
	if formFields.fileName != "" {
		if err := writer.WriteField("fileName", formFields.fileName); err != nil {
			return nil, err
		}
	}
	if formFields.fileExtension != "" {
		if err := writer.WriteField("fileExtension", formFields.fileExtension); err != nil {
			return nil, err
		}
	}
	if formFields.md5Hash != "" {
		if err := writer.WriteField("md5Hash", formFields.md5Hash); err != nil {
			return nil, err
		}
	}
	if formFields.chunkIndex != "" {
		if err := writer.WriteField("chunkIndex", formFields.chunkIndex); err != nil {
			return nil, err
		}
	}
	if formFields.totalChunks != "" {
		if err := writer.WriteField("totalChunks", formFields.totalChunks); err != nil {
			return nil, err
		}
	}
	if len(formFields.chunkContent) > 0 {
		part, err := writer.CreateFormFile("chunk", formFields.fileName)
		if err != nil {
			return nil, err
		}
		if _, err := part.Write(formFields.chunkContent); err != nil {
			return nil, err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "/upload", &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func createRandomChunks(id string, n int, path string) (string, error) {
	cfg := config.LoadConfig()

	chunksDir := filepath.Join(path, cfg.ChunksDir, id)
	if err := os.MkdirAll(chunksDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("Error while creating chunks directory %v", err)
	}

	byteSize := 1024*1024*5
	chunkContent := make([]byte, byteSize)

	hash := md5.New()

	for chunkIndex := 0; chunkIndex < n; chunkIndex++ {
		if _,err := rand.Read(chunkContent); err != nil {
			return "", fmt.Errorf("Received unexpected error while creating chunk content: %v", err)
		}
	
		chunkFilePath := filepath.Join(chunksDir, fmt.Sprintf("chunk_%d", chunkIndex))
		out, err := os.Create(chunkFilePath)
		if err != nil {
			return "", fmt.Errorf("Error while creating chunks file %v", err)
		}
		defer out.Close()
		if _, err := out.Write(chunkContent); err != nil {
			return "", fmt.Errorf("Error writing data to chunk file: %v", err)
        }

		if _, err := hash.Write(chunkContent); err != nil {
			return "", fmt.Errorf("error updating md5 hash: %v", err)
		}
	}

	finalHash := hex.EncodeToString(hash.Sum(nil))
	return finalHash, nil
}

func pathExists(filePath string) bool {
    _, err := os.Stat(filePath)
    if err != nil {
        if os.IsNotExist(err) {
            return false
        }
        fmt.Println("Error checking file:", err)
        return false
    }
    return true
}

// --------------------------------------
// 		  Suite Setup - Cleanup
// --------------------------------------
func TestMain(m *testing.M) {
	cfg := config.LoadConfig()
	if err := os.MkdirAll(cfg.SharingDir, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create sharing directory %q: %v\n", cfg.SharingDir, err)
		os.Exit(1)
	}
	if err := os.MkdirAll(cfg.UploadDir, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create upload directory %q: %v\n", cfg.UploadDir, err)
		os.Exit(1)
	}
	if err := os.MkdirAll(cfg.ChunksDir, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create upload directory %q: %v\n", cfg.ChunksDir, err)
		os.Exit(1)
	}
	if err := os.MkdirAll("secrets", os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create upload directory %q: %v\n", "secrets", err)
		os.Exit(1)
	}

	exitCode := m.Run()

	if err := os.RemoveAll(cfg.SharingDir); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to remove sharing directory %q: %v\n", cfg.SharingDir, err)
	}
	if err := os.RemoveAll(cfg.UploadDir); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to remove upload directory %q: %v\n", cfg.UploadDir, err)
	}
	if err := os.RemoveAll(cfg.ChunksDir); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to remove sharing directory %q: %v\n", cfg.ChunksDir, err)
	}
	if err := os.RemoveAll("secrets"); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to remove upload directory %q: %v\n", "secrets", err)
	}

	os.Exit(exitCode)
}

// --------------------------------------
// 			 Parse Form Tests
// --------------------------------------
func TestFormMissingFields(t *testing.T) {
	tests := []struct {
		missingField	string
		modifyForm		func(*FormFields)
		expectedErr		string
	} {
		{
			"fileId",
			func(f *FormFields) {
				f.fileId = ""
			},
			"fileId is required",
		},
		{
			"fileName",
			func(f *FormFields) {
				f.fileName = ""
			},
			"fileName is required",
		},
		{
			"fileExtension",
			func(f *FormFields) {
				f.fileExtension = ""
			},
			"fileExtension is required",
		},
		{
			"md5Hash",
			func(f *FormFields) {
				f.md5Hash = ""
			},
			"md5Hash is required",
		},
		{
			"chunkIndex",
			func(f *FormFields) {
				f.chunkIndex = ""
			},
			"chunkIndex is required",
		},
		{
			"totalChunks",
			func(f *FormFields) {
				f.totalChunks = ""
			},
			"totalChunks is required",
		},
		{
			"chunk",
			func(f *FormFields) {
				f.chunkContent = nil
			},
			"chunk file is required",
		},
	}

	for _, tt := range tests {
		t.Run("MissingField_"+tt.missingField, func (t *testing.T) {
			byteSize := 100
			form := FormFields{
				fileId:			uuid.New().String(),
				fileName:		"someFileName",
				fileExtension:	".txt",
				md5Hash:		"9c768e67e63a8e1762f2799cde1d912e",
				chunkIndex:		"0",
				totalChunks:	"1",
				chunkContent: 	make([]byte, byteSize),
			}
			tt.modifyForm(&form)
			req, err := createMultipartForm(form)
			if err != nil {
				t.Errorf("Received unexpected error when creating multipart form %v", err)
			}
			recorder := httptest.NewRecorder()

			_, _, err = ParseForm(recorder, req)
			if err == nil || err.Error() != tt.expectedErr {
				t.Errorf("Expected error %s for missing field %s, received %v", tt.expectedErr, tt.missingField, err)
			}
		})
	}
}

func TestUnsafeExtensions(t *testing.T) {
	tests := []struct {
		unsafeExtension		string
		expectedErr			string
	} {
		{".exe", "invalid file extension: .exe"},
		{".bat", "invalid file extension: .bat"},
		{".cmd", "invalid file extension: .cmd"},
		{".msi", "invalid file extension: .msi"},
		{".com", "invalid file extension: .com"},
		{".vbs", "invalid file extension: .vbs"},
		{".wsf", "invalid file extension: .wsf"},
		{".pif", "invalid file extension: .pif"},
		{".scr", "invalid file extension: .scr"},
		{".apk", "invalid file extension: .apk"},
		{".zip", "invalid file extension: .zip"},
		{".rar", "invalid file extension: .rar"},
		{".tar", "invalid file extension: .tar"},
		{".gz", "invalid file extension: .gz"},
		{".7z", "invalid file extension: .7z"},
		{".iso", "invalid file extension: .iso"},
		{".tar.gz", "invalid file extension: .tar.gz"},
		{".torrent", "invalid file extension: .torrent"},
		{".dll", "invalid file extension: .dll"},
		{".lib", "invalid file extension: .lib"},
		{".sys", "invalid file extension: .sys"},
		{".vbe", "invalid file extension: .vbe"},
		{".jse", "invalid file extension: .jse"},
		{".chm", "invalid file extension: .chm"},
		{".mswmm", "invalid file extension: .mswmm"},
		{".reg", "invalid file extension: .reg"},
		{".php", "invalid file extension: .php"},
		{".asp", "invalid file extension: .asp"},
		{".jsp", "invalid file extension: .jsp"},
		{".cgi", "invalid file extension: .cgi"},
		{".pl", "invalid file extension: .pl"},
		{".sh", "invalid file extension: .sh"},
	}
	for _, tt := range tests {
		t.Run("Unsafe_extension_"+tt.unsafeExtension, func(t *testing.T) {
			byteSize := 100
			form := FormFields{
				fileId:			uuid.New().String(),
				fileName:		"someFileName",
				fileExtension:	tt.unsafeExtension,
				md5Hash:		"9c768e67e63a8e1762f2799cde1d912e",
				chunkIndex:		"0",
				totalChunks:	"1",
				chunkContent: 	make([]byte, byteSize),
			}
			req, err := createMultipartForm(form)
			if err != nil {
				t.Errorf("Received unexpected error when creating multipart form %v", err)
			}
			recorder := httptest.NewRecorder()

			_, _, err = ParseForm(recorder, req)
			if err == nil || err.Error() != tt.expectedErr {
				t.Errorf("Expected error %s for unsafe file extension %s, received %s", tt.expectedErr, tt.unsafeExtension, err)
			}
		})
	}

}

func TestTooLargeFormPart(t *testing.T) {
	byteSize := 1024*1024*10 // 10MB - should fail
	form := FormFields{
		fileId:			uuid.New().String(),
		fileName:		"someFileName",
		fileExtension:	".txt",
		md5Hash:		"9c768e67e63a8e1762f2799cde1d912e",
		chunkIndex:		"0",
		totalChunks:	"1",
		chunkContent: 	make([]byte, byteSize),
	}
	req, err := createMultipartForm(form)
	if err != nil {
		t.Errorf("Received unexpected error when creating multipart form %v", err)
	}
	recorder := httptest.NewRecorder()

	_, _, err = ParseForm(recorder, req)
	if err == nil || !strings.Contains(err.Error(), "unable to parse form:") {
		t.Errorf("Expected error `unable to parse form:` for too large file, received %v", err)
	}
}

// --------------------------------------
// 		  Chunk Assembling tests
// --------------------------------------

func TestChunkAssemble(t *testing.T) {
	cfg := config.LoadConfig()

	t.Run("Chunk_Assemble_Invalid_Md5_Hash", func (t *testing.T) {
		id := uuid.New().String()

		if _, err := createRandomChunks(id, 5, cfg.UploadDir); err != nil {
			t.Error(err)
		}
		meta := ChunkMeta{
			FileId: 		id,
			FileName:      	"someName",
			FileExtension: 	".txt",
			MD5Hash:       	"deadbeef",
			ChunkIndex:   	4,
			TotalChunks:  	5,
		}
		jm := job.NewJobManager(30 * time.Minute)
		ChunkAssemble(meta, jm, cfg.UploadDir)
		
		fileName := meta.FileName + meta.FileExtension
		finalFilePath := filepath.Join(cfg.UploadDir, fileName)
		if pathExists(finalFilePath) {
			t.Error("Final file was created eventhough the md5 hash was wrong")
		}
	})

	t.Run("Chunk_Assemble_Success", func (t *testing.T) {
		id := uuid.New().String()
		hash, err := createRandomChunks(id, 5, cfg.UploadDir)
		if err != nil {
			t.Error(err)
		}
		meta := ChunkMeta{
			FileId: 		id,
			FileName:      	"someName",
			FileExtension: 	".txt",
			MD5Hash:       	hash,
			ChunkIndex:   	4,
			TotalChunks:  	5,
		}
		jm := job.NewJobManager(30 * time.Minute)
		ChunkAssemble(meta, jm, cfg.UploadDir)
		
		fileName := meta.FileName + meta.FileExtension
		finalFilePath := filepath.Join(cfg.UploadDir, fileName)
		if !pathExists(finalFilePath) {
			t.Error("Chunk assembly failed: Final File wasn't created")
		}
		chunksDir := filepath.Join(cfg.ChunksDir, id)
		if pathExists(chunksDir) {
			t.Error("Chunk assembly failed: Chunks Folder wasn't deleted")
		}
	})
}
// --------------------------------------
// 		  Authorization Tests
// --------------------------------------
func TestUploadHandlerAuth(t *testing.T) {
	// Create both folders to test
	folder1 := uuid.New().String()
	folder2 := uuid.New().String()
	if err := os.MkdirAll(folder1, os.ModePerm); err != nil {
		t.Fatalf("Received unexpected error when creating folder: %v", err)
	}
	if err := os.MkdirAll(folder2, os.ModePerm); err != nil {
		t.Fatalf("Received unexpected error when creating folder: %v", err)
	}
	form := FormFields{
		fileId:			uuid.New().String(),
		fileName:		"someFileName",
		fileExtension:	".txt",
		md5Hash:		"6d0bb00954ceb7fbee436bb55a8397a9",
		chunkIndex:		"0",
		totalChunks:	"1",
		chunkContent: 	make([]byte, 100),
	}
	expectedBody := "Forbidden: insufficient permissions"
	t.Run("Upload_Handler_Incorrect_Folder_Auth", func (t *testing.T) {
		claims := jwt.MapClaims{
			"user_id":   "someRandomUser",
			"folder_id": folder1,
			"access":    "w",
			"exp":       time.Now().Add(5 * time.Hour).Unix(),
		}
		ctx := context.WithValue(context.Background(), auth.ClaimsContextKey, claims)

		req, err := createMultipartForm(form)
		if err != nil {
			t.Fatalf("Received unexpected error when creating multipart form %v", err)
		}
		rr := httptest.NewRecorder()
		req = req.WithContext(ctx)
		jm := job.NewJobManager(30 * time.Minute)

		UploadHandler(rr, req, jm, folder2)

		if rr.Code != http.StatusForbidden {
			t.Errorf("expected status 403 Forbidden; got %d", rr.Code)
		}
		if strings.TrimSpace(rr.Body.String()) != expectedBody {
			t.Errorf("expected response body %q; got %q", expectedBody, strings.TrimSpace(rr.Body.String()))
		}
	})
	t.Run("Upload_Handler_Incorrect_Auth_Access", func (t *testing.T) {
		claims := jwt.MapClaims{
			"user_id":   "someRandomUser",
			"folder_id": folder1,
			"access":    "r",
			"exp":       time.Now().Add(5 * time.Hour).Unix(),
		}
		ctx := context.WithValue(context.Background(), auth.ClaimsContextKey, claims)

		req, err := createMultipartForm(form)
		if err != nil {
			t.Fatalf("Received unexpected error when creating multipart form %v", err)
		}
		rr := httptest.NewRecorder()
		req = req.WithContext(ctx)
		jm := job.NewJobManager(30 * time.Minute)

		UploadHandler(rr, req, jm, folder1)

		if rr.Code != http.StatusForbidden {
			t.Errorf("expected status 403 Forbidden; got %d", rr.Code)
		}
		if strings.TrimSpace(rr.Body.String()) != expectedBody {
			t.Errorf("expected response body %q; got %q", expectedBody, strings.TrimSpace(rr.Body.String()))
		}
	})
	t.Run("Successful_Authorization", func (t *testing.T) {
		claims := jwt.MapClaims{
			"user_id":   "someRandomUser",
			"folder_id": folder1,
			"access":    "w",
			"exp":       time.Now().Add(5 * time.Hour).Unix(),
		}
		ctx := context.WithValue(context.Background(), auth.ClaimsContextKey, claims)

		req, err := createMultipartForm(form)
		if err != nil {
			t.Fatalf("Received unexpected error when creating multipart form %v", err)
		}
		rr := httptest.NewRecorder()
		req = req.WithContext(ctx)
		jm := job.NewJobManager(30 * time.Minute)

		UploadHandler(rr, req, jm, folder1)

		if rr.Code == http.StatusForbidden {
			t.Errorf("didn't expect status 403 forbidden; got %d", rr.Code)
		}
	})

	time.Sleep(100 * time.Millisecond)

	if err := os.RemoveAll(folder1); err != nil {
		t.Fatalf("Received unexpected error when removing folder %s: %v", folder1, err)
	}
	if err := os.RemoveAll(folder2); err != nil {
		t.Fatalf("Received unexpected error when removing folder %s: %v", folder2, err)
	}
}

func TestUploadHandlerSuccess(t *testing.T) {
	cfg := config.LoadConfig()
	folder := uuid.New().String()
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		t.Fatalf("Received unexpected error when creating folder: %v", err)
	}
	claims := jwt.MapClaims{
		"user_id":   "someRandomUser",
		"folder_id": "/",
		"access":    "w",
		"exp":       time.Now().Add(5 * time.Hour).Unix(),
	}
	ctx := context.WithValue(context.Background(), auth.ClaimsContextKey, claims)

	form := FormFields{
		fileId:			uuid.New().String(),
		fileName:		"someFileName",
		fileExtension:	".txt",
		md5Hash:		"6d0bb00954ceb7fbee436bb55a8397a9",
		chunkIndex:		"0",
		totalChunks:	"1",
		chunkContent: 	make([]byte, 100),
	}
	req, err := createMultipartForm(form)
	if err != nil {
		t.Fatalf("Received unexpected error when creating multipart form %v", err)
	}
	rr := httptest.NewRecorder()
	req = req.WithContext(ctx)
	jm := job.NewJobManager(30 * time.Minute)

	UploadHandler(rr, req, jm, cfg.UploadDir)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200 OK; got %d", rr.Code)
	}

	
	time.Sleep(100 * time.Millisecond)

	fullFilePath := filepath.Join(cfg.UploadDir, "someFileName.txt")
	if !pathExists(fullFilePath) {
		t.Errorf("Final file wasn't created in: %s ", fullFilePath)
	}

	if err := os.RemoveAll(folder); err != nil {
		t.Fatalf("Received unexpected error when removing folder %s: %v", folder, err)
	}
}