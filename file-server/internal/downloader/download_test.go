package downloader

import (
	"context"
	"strings"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"

	"file-server/internal/auth"
	"file-server/internal/job"
	"file-server/config"
)

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


func TestDownloaderMissingParameters(t *testing.T) {
	t.Run("Missing_Parameters_Folder_id", func (t *testing.T) {
		path := "/login"
		queryParams := url.Values{}
		queryParams.Add("folder_id", "somefolderId")

		req := httptest.NewRequest(http.MethodGet, path+"?"+queryParams.Encode(), nil)
		rr := httptest.NewRecorder()

		jm := job.NewJobManager(30 * time.Minute)

		DownloadHandler(rr, req, jm)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 Bad Request, received %d", rr.Code)
		}
		if strings.TrimSpace(rr.Body.String()) != "Missing file parameter" {
			t.Errorf("Received invalid response body: %s", strings.TrimSpace(rr.Body.String()))
		}
	})

	t.Run("Missing_Parameters_File", func (t *testing.T) {
		path := "/login"
		queryParams := url.Values{}
		queryParams.Add("file", "someFileName")

		req := httptest.NewRequest(http.MethodGet, path+"?"+queryParams.Encode(), nil)
		rr := httptest.NewRecorder()

		jm := job.NewJobManager(30 * time.Minute)

		DownloadHandler(rr, req, jm)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 Bad Request, received %d", rr.Code)
		}
		if strings.TrimSpace(rr.Body.String()) != "Missing folder_id parameter" {
			t.Errorf("Received invalid response body: %s", strings.TrimSpace(rr.Body.String()))
		}
	})
}

func TestDownloaderNonExistentFile(t *testing.T) {
	folderId := uuid.New().String()
	path := "/login"
	queryParams := url.Values{}
	queryParams.Add("folder_id", folderId)
	queryParams.Add("file", "someNonExistentFileName")

	claims := jwt.MapClaims{
		"user_id":   "someRandomUser",
		"folder_id": folderId,
		"access":    "r",
		"exp":       time.Now().Add(5 * time.Hour).Unix(),
	}
	ctx := context.WithValue(context.Background(), auth.ClaimsContextKey, claims)

	req := httptest.NewRequest(http.MethodGet, path+"?"+queryParams.Encode(), nil)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	jm := job.NewJobManager(30 * time.Minute)

	DownloadHandler(rr, req, jm)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 Not Found, received %d", rr.Code)
	}
	if strings.TrimSpace(rr.Body.String()) != "File not found" {
		t.Errorf("Received invalid response body: %s", strings.TrimSpace(rr.Body.String()))
	}
}

func TestDownloaderInvalidAuth(t *testing.T) {
	t.Run("Invalid_Auth_Different_Folder", func (t *testing.T) {
		folder1 := uuid.New().String()
		folder2 := uuid.New().String()
		claims := jwt.MapClaims{
			"user_id":   "someRandomUser",
			"folder_id": folder1,
			"access":    "r",
			"exp":       time.Now().Add(5 * time.Hour).Unix(),
		}
		ctx := context.WithValue(context.Background(), auth.ClaimsContextKey, claims)

		path := "/login"
		queryParams := url.Values{}
		queryParams.Add("folder_id", folder2)
		queryParams.Add("file", "someNonExistentFileName")

		req := httptest.NewRequest(http.MethodGet, path+"?"+queryParams.Encode(), nil)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		jm := job.NewJobManager(30 * time.Minute)

		DownloadHandler(rr, req, jm)

		if rr.Code != http.StatusForbidden {
			t.Errorf("Expected status 403 Forbidden, received %d", rr.Code)
		}
		if strings.TrimSpace(rr.Body.String()) != "Forbidden: insufficient permissions" {
			t.Errorf("Received invalid response body: %s", strings.TrimSpace(rr.Body.String()))
		}
	})

	t.Run("Invalid_Auth_Different_Access", func (t *testing.T) {
		folder := uuid.New().String()
		claims := jwt.MapClaims{
			"user_id":   "someRandomUser",
			"folder_id": folder,
			"access":    "w", //Needs r access
			"exp":       time.Now().Add(5 * time.Hour).Unix(),
		}
		ctx := context.WithValue(context.Background(), auth.ClaimsContextKey, claims)

		path := "/login"
		queryParams := url.Values{}
		queryParams.Add("folder_id", folder)
		queryParams.Add("file", "someNonExistentFileName")

		req := httptest.NewRequest(http.MethodGet, path+"?"+queryParams.Encode(), nil)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		jm := job.NewJobManager(30 * time.Minute)

		DownloadHandler(rr, req, jm)

		if rr.Code != http.StatusForbidden {
			t.Errorf("Expected status 403 Forbidden, received %d", rr.Code)
		}
		if strings.TrimSpace(rr.Body.String()) != "Forbidden: insufficient permissions" {
			t.Errorf("Received invalid response body: %s", strings.TrimSpace(rr.Body.String()))
		}
	})
}

func TestDownloaderSuccess(t *testing.T) {
	folder := uuid.New().String()
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		t.Fatalf("Received unexpected error when creating folder: %v", err)
	}
	fileName := "someRandomFileName.txt"
	fullPath := filepath.Join(folder, fileName)

	out, err := os.Create(fullPath)
	if err != nil {
		t.Fatalf("Received unexpected error when creating file: %v", err)
	}

	byteSize := 1024*1024*100 // 100MB
	chunkContent := make([]byte, byteSize)
	n, err := out.Write(chunkContent)
	if err != nil {
		t.Fatalf("Error writing to file: %v", err)
	}
	if n != len(chunkContent) {
		t.Fatalf("Expected to write %d bytes but wrote %d", len(chunkContent), n)
	}
	out.Close()

	claims := jwt.MapClaims{
		"user_id":   "someRandomUser",
		"folder_id": folder,
		"access":    "r",
		"exp":       time.Now().Add(5 * time.Hour).Unix(),
	}
	ctx := context.WithValue(context.Background(), auth.ClaimsContextKey, claims)
	
	path := "/login"
	queryParams := url.Values{}
	queryParams.Add("folder_id", folder)
	queryParams.Add("file", fileName)

	req := httptest.NewRequest(http.MethodGet, path+"?"+queryParams.Encode(), nil)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	jm := job.NewJobManager(30 * time.Minute)

	DownloadHandler(rr, req, jm)

	if err := os.RemoveAll(folder); err != nil {
		t.Fatalf("Received unexpected error when deleting folder: %v", err)
	}
}