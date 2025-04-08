package downloader

import (
	"os"
	"path/filepath"
	"strconv"
	"net/http"
	
	"file-server/config"
	"file-server/internal/auth"
	"file-server/internal/job"

	
	"github.com/golang-jwt/jwt/v5"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request, jm *job.JobManager) {
	cfg := config.LoadConfig()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "Missing file parameter", http.StatusBadRequest)
		return
	}

	folderId := r.URL.Query().Get("folder_id")
	if folderId == "" {
		http.Error(w, "Missing folder_id parameter", http.StatusBadRequest)
		return
	}

	// Check if user has write access to the folder
	claimsRaw := r.Context().Value(auth.ClaimsContextKey)
	claims, ok := claimsRaw.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	canAccess, err := auth.HasAccess(claims, folderId, "r")
	if err != nil || !canAccess {
		http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
		return
	}

	filePath := filepath.Join(cfg.SharingDir, folderId, fileName)
	f, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
	}
	
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filepath.Base(fileName)+"\"")
	w.Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))

	http.ServeContent(w, r, fileName, fi.ModTime(), f)
}

func DownloadAvailableHandler(w http.ResponseWriter, r *http.Request, jm *job.JobManager) {
	cfg := config.LoadConfig()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "Missing file parameter", http.StatusBadRequest)
		return
	}

	folderId := r.URL.Query().Get("folder_id")
	if folderId == "" {
		http.Error(w, "Missing folder_id parameter", http.StatusBadRequest)
		return
	}

	// Check if user has write access to the folder
	claimsRaw := r.Context().Value(auth.ClaimsContextKey)
	claims, ok := claimsRaw.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	canAccess, err := auth.HasAccess(claims, folderId, "r")
	if err != nil || !canAccess {
		http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
		return
	}

	if (!jm.AcquireJob(fileName)) {
		http.Error(w, "File currently processing", http.StatusServiceUnavailable)
		return
	}
	jm.ReleaseJob(fileName)

	filePath := filepath.Join(cfg.SharingDir, folderId, fileName)
	f, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer f.Close()

	_, err = f.Stat()
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
	}
}