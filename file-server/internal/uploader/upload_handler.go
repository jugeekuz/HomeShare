package uploader

import (
	// "bufio"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"file-server/internal/auth"

	"github.com/google/uuid"	
	"github.com/golang-jwt/jwt/v5"
)
type FileMeta struct {
	FileId			string;
	FileName		string;
	FileExtension	string;
}
type File struct {
	File multipart.File
}

func ParseForm(w http.ResponseWriter, r *http.Request) (FileMeta, File, error) {
	const MAX_MBYTES = 5

	r.Body = http.MaxBytesReader(w, r.Body, (MAX_MBYTES<<20)+1024)

	if err := r.ParseMultipartForm(MAX_MBYTES << 20); err != nil {
		return FileMeta{}, File{}, fmt.Errorf("unable to parse form: %w", err)
	}

	if r.FormValue("fileId") == "" {
		return FileMeta{}, File{}, fmt.Errorf("fileId is required")
	}
	if r.FormValue("fileName") == "" {
		return FileMeta{}, File{}, fmt.Errorf("fileName is required")
	}
	if r.FormValue("fileExtension") == "" {
		return FileMeta{}, File{}, fmt.Errorf("fileExtension is required")
	}

	files, ok := r.MultipartForm.File["file"]
	if !ok || len(files) == 0 {
		return FileMeta{}, File{}, fmt.Errorf("file is required")
	}
	if files[0].Size == 0 {
		return FileMeta{}, File{}, fmt.Errorf("file is empty")
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		return FileMeta{}, File{}, fmt.Errorf("error while reading chunk: %w", err)
	}

	meta := FileMeta{
		FileId:        r.FormValue("fileId"),
		FileName:      r.FormValue("fileName"),
		FileExtension: r.FormValue("fileExtension"),
	}

	fileNameRegex := regexp.MustCompile(`^[a-zA-Z0-9._ -]+$`)
	if !fileNameRegex.MatchString(meta.FileName) {
		return FileMeta{}, File{}, fmt.Errorf("invalid file name format: %s", meta.FileName)
	}

	fileExtensionRegex := regexp.MustCompile(`^\.(jpe?g|png|gif|bmp|tiff?|webp|mp4|mkv|mov|avi|flv|wmv|txt)$`)
	if !fileExtensionRegex.MatchString(strings.ToLower(meta.FileExtension)) {
		return FileMeta{}, File{}, fmt.Errorf("invalid file extension: %s", meta.FileExtension)
	}

	if _, err := uuid.Parse(meta.FileId); err != nil {
		return FileMeta{}, File{}, fmt.Errorf("invalid UUID format: %w", err)
	}

	chunk := File{
		File: file,
	}

	return meta, chunk, nil
}

func UploadHandler(w http.ResponseWriter, r *http.Request, filePath string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user has write access to the folder
	claimsRaw := r.Context().Value(auth.ClaimsContextKey)
	claims, ok := claimsRaw.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	canAccess, err := auth.HasAccess(claims, "/", "w")
	if err != nil || !canAccess {
		http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
		return
	}

	meta, file, err := ParseForm(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.File.Close()

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		http.Error(w, "Unable to create upload directory", http.StatusInternalServerError)
		return
	}
	finalFilePath := filepath.Join(filePath, meta.FileName+meta.FileExtension)
	finalFilePath = getUniqueFileName(finalFilePath)

	destFile, err := os.Create(finalFilePath)
	if err != nil {
		http.Error(w, "Unable to create destination file", http.StatusInternalServerError)
		return
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, file.File); err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}