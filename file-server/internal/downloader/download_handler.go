package downloader

import (
	"os"
	"path/filepath"
	"strconv"
	"net/http"
	
	"file-server/config"
	"file-server/internal/job"
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