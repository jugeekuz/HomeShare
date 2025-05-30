package sharing

import (
	"bytes"
	"context"
	"database/sql"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
	
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"file-server/config"
	"file-server/internal/auth"
	"file-server/internal/helpers"
	"file-server/internal/job"
)

type FormFields struct {
	fileId        string
	fileName      string
	fileExtension string
	md5Hash       string
	chunkIndex    string
	totalChunks   string
	chunkContent  []byte
}

// --------------------------------------
//
//	Helper Functions
//
// --------------------------------------

func initMockDb() (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	return db, mock, nil
}

func createMultipartForm(url string, formFields FormFields) (*http.Request, error) {
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

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

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

// Sharing Handler Tests
func createSharingReq(claimFolderId string, folderName string, access string, expirationDate string, otpPass string) (*http.Request, error) {
	cfg := config.LoadConfig()
	url := "/share"

	claims := jwt.MapClaims{
		"user_id":   "someRandomUser",
		"folder_id": claimFolderId,
		"access":    access,
		"exp":       time.Now().Add(cfg.Secrets.Jwt.AccessExpiryDuration).Unix(),
	}
	ctx := context.WithValue(context.Background(), auth.ClaimsContextKey, claims)

	creds := SharingDetails{
		FolderName: folderName,
		OtpPass: otpPass,
		Access: access,
		ExpirationDate: expirationDate,
	}
	body, err := json.Marshal(creds)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal credentials: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	req = req.WithContext(ctx)

	return req, nil
}

func validateSharingToken(tokenString string) (auth.TokenParameters, error) {
	cfg := config.LoadConfig()

	var response auth.TokenParameters

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.Secrets.Jwt.JwtSecret), nil
	})
	if err != nil {
		return response, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if userId, ok := claims["user_id"].(string); ok {
			response.UserId = userId
		} else {
			return response, fmt.Errorf("user_id claim missing or invalid")
		}

		if expiryDuration, ok := claims["exp"].(float64); ok {
			response.ExpiryDuration = time.Duration(expiryDuration) * time.Second
		} else {
			return response, fmt.Errorf("exp claim missing or invalid")
		}

		if folderId, ok := claims["folder_id"].(string); ok {
			response.FolderId = folderId
		} else {
			return response, fmt.Errorf("folder_id claim missing or invalid")
		}

		if access, ok := claims["access"].(string); ok {
			response.Access = access
		} else {
			return response, fmt.Errorf("access claim missing or invalid")
		}
		return response, nil
	}

	return response, fmt.Errorf("invalid token")
}

// Test the auth functionality of the caller of the Sharing Handler
func TestSharingAuth(t *testing.T) {
	t.Run("Test_Sharing_Auth_No_Root_Access", func(t *testing.T) {
		rr := httptest.NewRecorder()

		now := time.Now().UTC()
		req, err := createSharingReq("someFolderId", "someFolderName", "rw", now.Add(48 * time.Hour).Format(time.RFC3339), "123456")
		if err != nil {
			t.Fatalf("Received unexpected error when creating request: %v", err)
		}

		db, _, err := initMockDb()
		if err != nil {
			t.Fatalf("Received unexpected error when initializing mock db: %v", err)
		}
		defer db.Close()

		SharingHandler(rr, req, db, "someSalt", "someLink")

		if rr.Code != http.StatusForbidden {
			t.Errorf("expected status 403 Forbidden, got: %d", rr.Code)
		}
		if strings.TrimSpace(rr.Body.String()) != "Forbidden: insufficient permissions" {
			t.Errorf("expected error `Forbidden: insufficient permissions`, received: %s", rr.Body.String())
		}
	})

	t.Run("Test_Sharing_Auth_No_RW_Access", func(t *testing.T) {
		rr := httptest.NewRecorder()

		now := time.Now().UTC()
		req, err := createSharingReq("someFolderId", "someFolderName", "w", now.Add(48 * time.Hour).Format(time.RFC3339), "123456")
		if err != nil {
			t.Fatalf("Received unexpected error when creating request: %v", err)
		}

		db, _, err := initMockDb()
		if err != nil {
			t.Fatalf("Received unexpected error when initializing mock db: %v", err)
		}
		defer db.Close()

		SharingHandler(rr, req, db, "someSalt", "someLink")

		if rr.Code != http.StatusForbidden {
			t.Errorf("expected status 403 Forbidden, got: %d", rr.Code)
		}
		if strings.TrimSpace(rr.Body.String()) != "Forbidden: insufficient permissions" {
			t.Errorf("expected error `Forbidden: insufficient permissions`, received: %s", rr.Body.String())
		}
	})

	t.Run("Test_Sharing_Auth_Success", func(t *testing.T) {
		rr := httptest.NewRecorder()
		now := time.Now().UTC()
		req, err := createSharingReq("/", "someFolderName", "rw", now.Add(48 * time.Hour).Format(time.RFC3339), "123456")
		if err != nil {
			t.Fatalf("Received unexpected error when creating request: %v", err)
		}

		db, _, err := initMockDb()
		if err != nil {
			t.Fatalf("Received unexpected error when initializing mock db: %v", err)
		}
		defer db.Close()

		SharingHandler(rr, req, db, "someSalt", "someLink")

		if rr.Code == http.StatusForbidden {
			t.Errorf("didn't expect status 403 Forbidden, got: %d", rr.Code)
		}
	})
}

// Test the permissions of the token returned
func TestCreateSharingSuccess(t *testing.T) {
	cfg := config.LoadConfig()

	rr := httptest.NewRecorder()

	expirationDate := time.Now().Add(48 * time.Hour).UTC().Truncate(time.Second)

	expiration := expirationDate.Format(time.RFC3339)
	otpPass := "123456"
	linkUrl := uuid.New().String()
	expiryDuration := expirationDate.Sub(time.Now().UTC().Truncate(time.Second))
	folderId := helpers.GenerateFolderName(expiryDuration, linkUrl)
	folderName := "someFolderName"
	access := "rw"
	salt, err := helpers.GenerateRandomSalt()
	if err != nil {
		t.Errorf("Received unexpected error when generating random salt: %v", err)
	}
	hashedOtp := helpers.HashPassword(otpPass, salt)

	req, err := createSharingReq("/", folderName, "rw", expiration, otpPass)
	if err != nil {
		t.Fatalf("Received unexpected error when creating request: %v", err)
	}

	db, mock, err := initMockDb()
	if err != nil {
		t.Fatalf("Received unexpected error when initializing mock db: %v", err)
	}
	defer db.Close()
	
	mock.ExpectExec(`INSERT INTO sharing_users \(link_url, folder_id, folder_name, salt, otp_hash, access, expiration\)[\s\n]*VALUES[\s\n]*\(\$1, \$2, \$3, \$4, \$5, \$6, \$7\)[\s\n]*ON CONFLICT[\s\n]*\(link_url\)[\s\n]*DO UPDATE[\s\n]*SET link_url = EXCLUDED.link_url[\s\n]*RETURNING link_url, folder_id, folder_name, salt, otp_hash, access, expiration`).
		WithArgs(
			linkUrl, folderId, folderName, salt, hashedOtp, access, expiration,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	SharingHandler(rr, req, db, salt, linkUrl)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200 OK, got: %d", rr.Code)
	}

	var sharingResponse SharingResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &sharingResponse); err != nil {
		t.Fatalf("Received unexpected error when parsing response body: %v", err)
	}

	if sharingResponse.LinkUrl != linkUrl {
		t.Errorf("Expected %s link url, got : %s", linkUrl, sharingResponse.LinkUrl)
	}
	if sharingResponse.FolderId != folderId {
		t.Errorf("Expected %s folder id, got : %s", linkUrl, sharingResponse.FolderId)
	}

	// Check if the folder was created
	fullFolderPath := filepath.Join(cfg.SharingDir, sharingResponse.FolderId)
	if _, err := os.Stat(fullFolderPath); err != nil {
		if os.IsNotExist(err) {
			t.Error("Sharing Folder wasn't created.")
		} else {
			t.Errorf("Received unexpected error when searching for sharing folder: %v", err)
		}
	}

}

// // Add Sharing Files Tests
func TestAddSharingFilesAuth(t *testing.T) {
	t.Run("Test_Add_Sharing_Auth_Wrong_Folder_Access", func(t *testing.T) {

		url := "/share-file"
		byteSize := 3 * 1024 * 1024
		jm := job.NewJobManager(30 * time.Minute)

		byteContent := make([]byte, byteSize)
		if _, err := rand.Read(byteContent); err != nil {
			t.Fatalf("error while reading into file: %v\n", err)
		}
		hash := md5.Sum(byteContent)

		form := FormFields{
			fileId:        uuid.New().String(),
			fileName:      "someFileName",
			fileExtension: ".txt",
			md5Hash:       hex.EncodeToString(hash[:]),
			chunkIndex:    "0",
			totalChunks:   "1",
			chunkContent:  byteContent,
		}
		rr := httptest.NewRecorder()
		req, err := createMultipartForm(url, form)
		if err != nil {
			t.Fatalf("Received unexpected error when creating multipart form: %v", err)
		}
		claims := jwt.MapClaims{
			"user_id":   "someRandomUser",
			"folder_id": "someFolderId",
			"access":    "rw",
			"exp":       time.Now().Add(30 * time.Minute).Unix(),
		}
		ctx := context.WithValue(context.Background(), auth.ClaimsContextKey, claims)
		req = req.WithContext(ctx)
		req.Header.Set("Folder-Id", "someOtherFolderId")

		AddSharingFilesHandler(rr, req, jm)

		if rr.Code != http.StatusForbidden {
			t.Errorf("Expected status 403 Forbidden, got: %d", rr.Code)
		}
		if strings.TrimSpace(rr.Body.String()) != "Forbidden: insufficient permissions" {
			t.Errorf("Expected `Forbidden: insufficient permissions`, got: %s", rr.Body.String())
		}

	})

	t.Run("Test_Add_Sharing_Auth_No_W_Access", func(t *testing.T) {
		url := "/share-file"
		byteSize := 3 * 1024 * 1024
		jm := job.NewJobManager(30 * time.Minute)

		byteContent := make([]byte, byteSize)
		if _, err := rand.Read(byteContent); err != nil {
			t.Fatalf("error while reading into file: %v\n", err)
		}
		hash := md5.Sum(byteContent)

		form := FormFields{
			fileId:        uuid.New().String(),
			fileName:      "someFileName",
			fileExtension: ".txt",
			md5Hash:       hex.EncodeToString(hash[:]),
			chunkIndex:    "0",
			totalChunks:   "1",
			chunkContent:  byteContent,
		}
		rr := httptest.NewRecorder()
		req, err := createMultipartForm(url, form)
		if err != nil {
			t.Fatalf("Received unexpected error when creating multipart form: %v", err)
		}
		claims := jwt.MapClaims{
			"user_id":   "someRandomUser",
			"folder_id": "someFolderId",
			"access":    "r",
			"exp":       time.Now().Add(30 * time.Minute).Unix(),
		}
		ctx := context.WithValue(context.Background(), auth.ClaimsContextKey, claims)
		req = req.WithContext(ctx)
		req.Header.Set("Folder-Id", "someFolderId")

		AddSharingFilesHandler(rr, req, jm)

		if rr.Code != http.StatusForbidden {
			t.Errorf("Expected status 403 Forbidden, got: %d", rr.Code)
		}
		if strings.TrimSpace(rr.Body.String()) != "Forbidden: insufficient permissions" {
			t.Errorf("Expected `Forbidden: insufficient permissions`, got: %s", rr.Body.String())
		}
	})
}

func TestAddSharingFolderNotExist(t *testing.T) {
	url := "/share-file"
	byteSize := 3 * 1024 * 1024
	jm := job.NewJobManager(30 * time.Minute)

	byteContent := make([]byte, byteSize)
	if _, err := rand.Read(byteContent); err != nil {
		t.Fatalf("error while reading into file: %v\n", err)
	}
	hash := md5.Sum(byteContent)

	form := FormFields{
		fileId:        uuid.New().String(),
		fileName:      "someFileName",
		fileExtension: ".txt",
		md5Hash:       hex.EncodeToString(hash[:]),
		chunkIndex:    "0",
		totalChunks:   "1",
		chunkContent:  byteContent,
	}
	rr := httptest.NewRecorder()
	req, err := createMultipartForm(url, form)
	if err != nil {
		t.Fatalf("Received unexpected error when creating multipart form: %v", err)
	}
	claims := jwt.MapClaims{
		"user_id":   "someRandomUser",
		"folder_id": "someFolderId",
		"access":    "rw",
		"exp":       time.Now().Add(30 * time.Minute).Unix(),
	}
	ctx := context.WithValue(context.Background(), auth.ClaimsContextKey, claims)
	req = req.WithContext(ctx)
	req.Header.Set("Folder-Id", "someFolderId")

	AddSharingFilesHandler(rr, req, jm)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 Not Found, got: %d", rr.Code)
	}
	if strings.TrimSpace(rr.Body.String()) != "Folder not found" {
		t.Errorf("Expected `Folder not found`, got: %s", rr.Body.String())
	}
}

func TestAddSharingFilesSuccess(t *testing.T) {
	cfg := config.LoadConfig()

	url := "/share-file"
	fileName := "someFileName"
	fileExt := ".txt"
	byteSize := 3 * 1024 * 1024
	jm := job.NewJobManager(30 * time.Minute)
	linkUrl := uuid.New().String()
	expirationDate := time.Now().Add(48 * time.Hour).UTC().Truncate(time.Second)
	expiryDuration := expirationDate.Sub(time.Now().UTC().Truncate(time.Second))
	sharingFolderId := helpers.GenerateFolderName(expiryDuration, linkUrl)

	// Create Folder
	finalSharingFolder := filepath.Join(cfg.SharingDir, sharingFolderId)
	if err := os.MkdirAll(finalSharingFolder, os.ModePerm); err != nil {
		t.Fatalf("Encounctered error while creating folder : %v", err)
		return
	}

	// Create file to upload
	byteContent := make([]byte, byteSize)
	if _, err := rand.Read(byteContent); err != nil {
		t.Fatalf("error while reading into file: %v\n", err)
	}
	hash := md5.Sum(byteContent)

	form := FormFields{
		fileId:        uuid.New().String(),
		fileName:      fileName,
		fileExtension: fileExt,
		md5Hash:       hex.EncodeToString(hash[:]),
		chunkIndex:    "0",
		totalChunks:   "1",
		chunkContent:  byteContent,
	}
	rr := httptest.NewRecorder()
	req, err := createMultipartForm(url, form)
	if err != nil {
		t.Fatalf("Received unexpected error when creating multipart form: %v", err)
	}

	claims := jwt.MapClaims{
		"user_id":   "someRandomUser",
		"folder_id": sharingFolderId,
		"access":    "rw",
		"exp":       time.Now().Add(30 * time.Minute).Unix(),
	}
	ctx := context.WithValue(context.Background(), auth.ClaimsContextKey, claims)
	req = req.WithContext(ctx)
	req.Header.Set("Folder-Id", sharingFolderId)

	AddSharingFilesHandler(rr, req, jm)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got: %d", rr.Code)
	}

	// Wait for the file to be assembled (test will return while chunk assemble works - server would have been live)
	time.Sleep(100 * time.Millisecond)

	fullFilePath := filepath.Join(cfg.SharingDir, sharingFolderId, fileName + fileExt)
	if _, err := os.Stat(fullFilePath); err != nil {
		if os.IsNotExist(err) {
			t.Error("File shared wasn't created.")
		} else {
			t.Errorf("Received unexpected error when searching for file: %v", err)
		}
	}

	// Wait for .zip file to be created
	time.Sleep(5 * time.Second)
	fullZipFilePath := filepath.Join(cfg.SharingDir, sharingFolderId, sharingFolderId + ".zip")
	if _, err := os.Stat(fullZipFilePath); err != nil {
		if os.IsNotExist(err) {
			t.Error("File shared wasn't created.")
		} else {
			t.Errorf("Received unexpected error when searching for file: %v", err)
		}
	}

}

// Get Sharing Files Tests
func TestGetSharingFilesAuth(t *testing.T) {
	urlPath := "/share-files"
	t.Run("Test_Get_Sharing_Auth_Wrong_Folder_Access", func(t *testing.T) {
		queryParams := url.Values{}
		queryParams.Add("folder_id", "someFolderId")
		claims := jwt.MapClaims{
			"user_id":   "someRandomUser",
			"folder_id": "someFolderId2",
			"access":    "rw",
			"exp":       time.Now().Add(30 * time.Minute).Unix(),
		}
		req := httptest.NewRequest(http.MethodGet, urlPath+"?"+queryParams.Encode(), nil)
		ctx := context.WithValue(context.Background(), auth.ClaimsContextKey, claims)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		GetSharingFilesHandler(rr, req)

		if rr.Code != http.StatusForbidden {
			t.Errorf("Expected status 403 Forbidden, got %d", rr.Code)
		}
		if strings.TrimSpace(rr.Body.String()) != "Forbidden: insufficient permissions" {
			t.Errorf("Expected `Forbidden: insufficient permissions` body, got: %s", rr.Body.String())
		}
	})

	t.Run("Test_Get_Sharing_Auth_No_R_Access", func(t *testing.T) {
		queryParams := url.Values{}
		queryParams.Add("folder_id", "someFolderId")
		claims := jwt.MapClaims{
			"user_id":   "someRandomUser",
			"folder_id": "someFolderId",
			"access":    "w",
			"exp":       time.Now().Add(30 * time.Minute).Unix(),
		}
		req := httptest.NewRequest(http.MethodGet, urlPath+"?"+queryParams.Encode(), nil)
		ctx := context.WithValue(context.Background(), auth.ClaimsContextKey, claims)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		GetSharingFilesHandler(rr, req)

		if rr.Code != http.StatusForbidden {
			t.Errorf("Expected status 403 Forbidden, got %d", rr.Code)
		}
		if strings.TrimSpace(rr.Body.String()) != "Forbidden: insufficient permissions" {
			t.Errorf("Expected `Forbidden: insufficient permissions` body, got: %s", rr.Body.String())
		}
	})

	t.Run("Test_Get_Sharing_Auth_Success", func(t *testing.T) {
		queryParams := url.Values{}
		queryParams.Add("folder_id", "someFolderId")
		claims := jwt.MapClaims{
			"user_id":   "someRandomUser",
			"folder_id": "someFolderId",
			"access":    "rw",
			"exp":       time.Now().Add(30 * time.Minute).Unix(),
		}
		req := httptest.NewRequest(http.MethodGet, urlPath+"?"+queryParams.Encode(), nil)
		ctx := context.WithValue(context.Background(), auth.ClaimsContextKey, claims)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		GetSharingFilesHandler(rr, req)

		if rr.Code == http.StatusForbidden {
			t.Errorf("Didn't expect status 403 Forbidden, got %d", rr.Code)
		}
		if strings.TrimSpace(rr.Body.String()) == "Forbidden: insufficient permissions" {
			t.Errorf("Didn't expect `Forbidden: insufficient permissions` body, got: %s", rr.Body.String())
		}
	})
}

func TestGetSharingFilesMissingParameters(t *testing.T) {
	urlPath := "/share-files"
	t.Run("Test_Get_Sharing_No_Folder_Id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, urlPath, nil)
		claims := jwt.MapClaims{
			"user_id":   "someRandomUser",
			"folder_id": "someFolderId",
			"access":    "rw",
			"exp":       time.Now().Add(30 * time.Minute).Unix(),
		}
		ctx := context.WithValue(context.Background(), auth.ClaimsContextKey, claims)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		GetSharingFilesHandler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 Bad Request, got %d", rr.Code)
		}
		if strings.TrimSpace(rr.Body.String()) != "Missing folder_id parameter" {
			t.Errorf("Expected `Missing folder_id parameter` body, got: %s", rr.Body.String())
		}

	})
}

func TestGetSharingFilesInvalidParameters(t *testing.T) {
	urlPath := "/share-files"
	t.Run("Test_Get_Sharing_Non_Existent_Folder", func(t *testing.T) {
		queryParams := url.Values{}
		queryParams.Add("folder_id", "someFolderId")
		claims := jwt.MapClaims{
			"user_id":   "someRandomUser",
			"folder_id": "someFolderId",
			"access":    "rw",
			"exp":       time.Now().Add(30 * time.Minute).Unix(),
		}
		req := httptest.NewRequest(http.MethodGet, urlPath+"?"+queryParams.Encode(), nil)
		ctx := context.WithValue(context.Background(), auth.ClaimsContextKey, claims)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		GetSharingFilesHandler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 Bad Request, got %d", rr.Code)
		}
		if strings.TrimSpace(rr.Body.String()) != "Folder does not exist" {
			t.Errorf("Expected `Folder does not exist` body, got: %s", rr.Body.String())
		}
	})
}

func TestGetSharingFilesSuccess(t *testing.T) {
	cfg := config.LoadConfig()

	urlPath := "/share-files"
	fileName := "someFileName"
	fileExt := ".txt"
	byteSize := 3 * 1024 * 1024
	linkUrl := uuid.New().String()
	expirationDate := time.Now().Add(48 * time.Hour).UTC().Truncate(time.Second)
	expiryDuration := expirationDate.Sub(time.Now().UTC().Truncate(time.Second))
	sharingFolderId := helpers.GenerateFolderName(expiryDuration, linkUrl)

	// Create Folder
	finalSharingFolder := filepath.Join(cfg.SharingDir, sharingFolderId)
	if err := os.MkdirAll(finalSharingFolder, os.ModePerm); err != nil {
		t.Fatalf("Encounctered error while creating folder : %v", err)
	}

	// Create file to save
	byteContent := make([]byte, byteSize)
	if _, err := rand.Read(byteContent); err != nil {
		t.Fatalf("error while reading into file: %v\n", err)
	}

	fullFilePath := filepath.Join(finalSharingFolder, fileName + fileExt)
	file, err := os.Create(fullFilePath)
	if err != nil {
		t.Fatalf("Encounctered error while creating file : %v", err)
	}

	n, err := file.Write(byteContent)
	if err != nil || n == 0 {
		t.Fatalf("Encounctered error while writing into file : %v", err)
	}

	claims := jwt.MapClaims{
		"user_id":   "someRandomUser",
		"folder_id": sharingFolderId,
		"access":    "rw",
		"exp":       time.Now().Add(30 * time.Minute).Unix(),
	}
	ctx := context.WithValue(context.Background(), auth.ClaimsContextKey, claims)

	// Now finally check
	queryParams := url.Values{}
	queryParams.Add("folder_id", sharingFolderId)
	req := httptest.NewRequest(http.MethodGet, urlPath+"?"+queryParams.Encode(), nil)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	time.Sleep(800 * time.Millisecond)

	GetSharingFilesHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", rr.Code)
	}

	// Check if files were returned
	var sharingFilesResponse SharingFilesResponse
	if err := json.NewDecoder(rr.Body).Decode(&sharingFilesResponse); err != nil {
		t.Fatalf("error unmarshalling sharing files response: %v", err)
	}

	if len(sharingFilesResponse.Files) != 1 {
		t.Errorf("Expected total sharing files of 1, got: %d", len(sharingFilesResponse.Files))
	}

	// Check if one of the files corresponds to "someFileName.txt"
	found := false
	for _, file := range sharingFilesResponse.Files {
		// Reconstruct the file name from name and extension)
		if file.FileName+file.FileExtension == fileName+fileExt {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("`%s` doesn't exist inside the sharing folder", fileName+fileExt)
	}
}
