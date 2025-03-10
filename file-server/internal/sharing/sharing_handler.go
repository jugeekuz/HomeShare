package sharing

import (
	"fmt"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"file-server/config"
	"file-server/internal/auth"
	"file-server/internal/helpers"

	"file-server/internal/job"
	"file-server/internal/uploader"

	"github.com/google/uuid"
)

type SharingDetails struct {
	ExpiryDuration	string `json:"expiry_duration"`
}

type SharingFileParameters struct {
	FolderId	string `json:"folder_id"`
}

type SharingResponse struct {
	RefreshToken 	string `json:"refresh_token"`
	FolderId 		string `json:"folder_id"`
}

func SharingHandler(w http.ResponseWriter, r *http.Request) {
	cfg := config.LoadConfig()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var sharingDetails SharingDetails
	if err := json.NewDecoder(r.Body).Decode(&sharingDetails); err != nil {
		http.Error(w, "Unable to parse sharing parameters", http.StatusBadRequest)
		return
	}
	expiryDuration, err := time.ParseDuration(sharingDetails.ExpiryDuration)
	if err != nil {
		http.Error(w, "Error while parsing duration", http.StatusInternalServerError)
		return
	}

	sharingFolderName := helpers.GenerateFolderName(expiryDuration)
	finalSharingFolder := filepath.Join(cfg.SharingDir ,sharingFolderName)
	if err := os.MkdirAll(finalSharingFolder, os.ModePerm); err != nil {
		http.Error(w, "Error while creating folder", http.StatusInternalServerError)
		return
	}

	refreshParams := &auth.TokenParameters{
		UserId:         uuid.New().String(),
		ExpiryDuration: expiryDuration,
		FolderId:       "/",
		Access:         "r",
	}
	_, refreshToken, err := auth.GenerateTokens(refreshParams, refreshParams)
	if err != nil {
		http.Error(w, "Error while generating tokens", http.StatusInternalServerError)
		return
	}

	var sharingResponse SharingResponse
	sharingResponse.RefreshToken = refreshToken
	sharingResponse.FolderId = sharingFolderName
	if err := json.NewEncoder(w).Encode(&sharingResponse); err != nil {
		http.Error(w, "Error while generating response", http.StatusInternalServerError)
		return
	}
}

func AddSharingFilesHandler(w http.ResponseWriter, r *http.Request, jm *job.JobManager) {
	cfg := config.LoadConfig()

	folderId := r.Header.Get("Folder-Id")
    if folderId == "" {
        http.Error(w, "Folder-Id header field is required", http.StatusBadRequest)
        return
    }
	folderDir := filepath.Join(cfg.SharingDir, folderId)

    if _, err := os.Stat(folderDir); err != nil {
        if os.IsNotExist(err) {
            http.Error(w, "Folder not found", http.StatusNotFound)
            return
        }
        http.Error(w, "Error checking folder: "+err.Error(), http.StatusInternalServerError)
        return
    }

	uploader.UploadHandler(w, r, folderDir)


	entries, err := os.ReadDir(folderDir)
	if err != nil {
		http.Error(w, "Error while reading directory", http.StatusInternalServerError)
		return
	}
	files := make([]string, len(entries))
	for index, entry := range entries {
		files[index] = entry.Name()
	}

	zipFileName := fmt.Sprintf("%s.zip", folderId)
	go helpers.CreateZip(folderDir, zipFileName, files, jm)
}
