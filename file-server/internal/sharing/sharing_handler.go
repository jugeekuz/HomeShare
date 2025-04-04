package sharing

import (
	"fmt"
	"encoding/json"
	"strings"
	"strconv"
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
	"github.com/golang-jwt/jwt/v5"
)

type SharingDetails struct {
	Access			string `json:"access"`
	FolderName		string `json:"folder_name"`
	ExpiryDuration	string `json:"expiry_duration"`
}

type SharingFileParameters struct {
	FolderId	string `json:"folder_id"`
}

type SharingResponse struct {
	RefreshToken 	string `json:"refresh_token"`
	FolderId 		string `json:"folder_id"`
}

type SharingFileItem struct {
	FileName 		string `json:"file_name"`
	FileExtension 	string `json:"file_extension"`
	FileSize		string `json:"file_size"`
}

type SharingFilesResponse struct {
	Files		[]SharingFileItem `json:"files"`
}

func SharingHandler(w http.ResponseWriter, r *http.Request) {
	cfg := config.LoadConfig()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is admin (ie has rw access to root)
	claimsRaw := r.Context().Value(auth.ClaimsContextKey)
	claims, ok := claimsRaw.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	canAccess, err := auth.HasAccess(claims, "/", "rw")
	if err != nil || !canAccess {
		http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
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
		FolderId:       sharingFolderName,
		Access:         sharingDetails.Access,
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

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cfg := config.LoadConfig()

	folderId := r.Header.Get("Folder-Id")
    if folderId == "" {
        http.Error(w, "Folder-Id header field is required", http.StatusBadRequest)
        return
    }
	// Check if user has write access to the folder
	claimsRaw := r.Context().Value(auth.ClaimsContextKey)
	claims, ok := claimsRaw.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}
	
	fullFolderIdPath := filepath.Join(cfg.SharingDir, folderId)

	canAccess, err := auth.HasAccess(claims, folderId, "w")
	if err != nil || !canAccess {
		http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
		return
	}

    if _, err := os.Stat(fullFolderIdPath); err != nil {
        if os.IsNotExist(err) {
            http.Error(w, "Folder not found", http.StatusNotFound)
            return
        }
        http.Error(w, "Error checking folder: "+err.Error(), http.StatusInternalServerError)
        return
    }

	uploader.UploadHandler(w, r, jm, fullFolderIdPath)

	entries, err := os.ReadDir(fullFolderIdPath)
	if err != nil {
		http.Error(w, "Error while reading directory", http.StatusInternalServerError)
		return
	}
	files := make([]string, len(entries))
	for index, entry := range entries {
		files[index] = entry.Name()
	}

	zipFileName := fmt.Sprintf("%s.zip", folderId)
	go helpers.CreateZip(fullFolderIdPath, zipFileName, files, jm)
}

func GetSharingFilesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cfg := config.LoadConfig()

	folderId := r.URL.Query().Get("folder_id")
	if folderId == "" {
		http.Error(w, "Missing folder_id parameter", http.StatusBadRequest)
		return
	}

	claimsRaw := r.Context().Value(auth.ClaimsContextKey)
	claims, ok := claimsRaw.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	// Check if user has read access to the folder
	canAccess, err := auth.HasAccess(claims, folderId, "r")
	if err != nil || !canAccess {
		http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
		return
	}

	folderPath := filepath.Join(cfg.SharingDir, folderId)
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		http.Error(w, "Folder does not exist", http.StatusBadRequest)
		return
	}

	var sharingFilesResponse SharingFilesResponse
	for _, entry := range entries {
		if !entry.IsDir() {
			fileName := entry.Name()
			ext := filepath.Ext(fileName)
			
			fileInfo, err := entry.Info()
			if err != nil {
				continue
			}

			nameWithoutExt := strings.TrimSuffix(fileName, ext)
			
			sharingFilesResponse.Files = append(sharingFilesResponse.Files, SharingFileItem{
				FileName:      nameWithoutExt,
				FileExtension: ext,
				FileSize:      strconv.FormatInt(fileInfo.Size(), 10),
			})
		}
	}
	
	json.NewEncoder(w).Encode(sharingFilesResponse)
}