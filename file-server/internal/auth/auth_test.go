package auth

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	"file-server/config"
	"file-server/internal/helpers"
)

// Mock Database somehow
func initMockDb() (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	return db, mock, nil
}

func validateToken(tokenString string, expectedFolder string, expectedAccess string) error {
	cfg := config.LoadConfig()

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.Secrets.Jwt.JwtSecret), nil
	})
	if err != nil || !token.Valid {
		return fmt.Errorf("invalid access token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("token claims are not of type MapClaims")
	}

	folder, ok := claims["folder_id"].(string)
	if !ok {
		return fmt.Errorf("folder claim missing or not a string")
	}
	if folder != expectedFolder {
		return fmt.Errorf("invalid folder claim: expected %s, got %s", expectedFolder, folder)
	}

	access, ok := claims["access"].(string)
	if !ok {
		return fmt.Errorf("access claim missing or not a string")
	}
	if access != expectedAccess {
		return fmt.Errorf("invalid access claim: expected %s, got %s", expectedAccess, access)
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return fmt.Errorf("expiration (exp) claim missing or not a number")
	}
	if time.Now().Unix() > int64(expFloat) {
		return fmt.Errorf("token is expired")
	}

	return nil
}

// --------------------------------------
//
//	Suite Setup - Cleanup
//
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

func TestLoginHandler(t *testing.T) {
	t.Run("Login_Handler_Incorrect_Credentials", func(t *testing.T) {
		// Initialize DB
		db, mock, err := initMockDb()
		if err != nil {
			t.Fatalf("Received unexpected error when initializing mock db: %v", err)
		}
		defer db.Close()
		rows := sqlmock.NewRows([]string{"username", "email", "salt", "password_hash", "folder", "access"}).
			AddRow("johndoe", "johndoe@example.com", "somesalt", "passwordsomesalt", "/", "rw")

		mock.ExpectQuery("SELECT username, email, salt, password_hash, folder, access FROM users WHERE username = \\$1").
			WithArgs("johndoe").
			WillReturnRows(rows)

		// Initialize w, rr
		creds := Credentials{
			Username: "johndoe",
			Password: "somepassword",
		}
		body, err := json.Marshal(creds)
		if err != nil {
			t.Fatalf("failed to marshal credentials: %v", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		LoginHandler(rr, req, db)

		if rr.Code != http.StatusForbidden {
			t.Errorf("Expected 403 Forbidden, got : %d", rr.Code)
		}

		if strings.TrimSpace(rr.Body.String()) != "Forbidden: invalid credentials" {
			t.Errorf("Expected 'Forbidden', got : %q", rr.Body.String())
		}
	})

	t.Run("Login_Handler_Success", func(t *testing.T) {
		creds := Credentials{
			Username: "johndoe",
			Password: "somepassword",
		}

		// Initialize DB
		db, mock, err := initMockDb()
		if err != nil {
			t.Fatalf("Received unexpected error when initializing mock db: %v", err)
		}
		defer db.Close()

		salt, err := helpers.GenerateRandomSalt()
		if err != nil {
			t.Fatalf("Encountered unexpected error while generating salt: %v", err)
		}
		hash := helpers.HashPassword(creds.Password, salt)

		rows := sqlmock.NewRows([]string{"username", "email", "salt", "password_hash", "folder", "access"}).
			AddRow("johndoe", "johndoe@example.com", salt, hash, "/", "rw")

		mock.ExpectQuery("SELECT username, email, salt, password_hash, folder, access FROM users WHERE username = \\$1").
			WithArgs("johndoe").
			WillReturnRows(rows)

		body, err := json.Marshal(creds)
		if err != nil {
			t.Fatalf("failed to marshal credentials: %v", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		LoginHandler(rr, req, db)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200 OK, got : %d", rr.Code)
		}

		// Check if access token, refresh token is set and valid.
		respData := TokenResponse{}
		if err := json.Unmarshal(rr.Body.Bytes(), &respData); err != nil {
			t.Fatalf("error unmarshalling response body: %v", err)
		}
		if respData.AccessToken == "" {
			t.Error("access token not found in response body")
		}

		if err := validateToken(respData.AccessToken, "/", "rw"); err != nil {
			t.Errorf("Unexpected error when validating access token: %v", err)
		}

		var refreshCookie *http.Cookie
		for _, cookie := range rr.Result().Cookies() {
			if cookie.Name == "refresh_token" {
				refreshCookie = cookie
				break
			}
		}
		if refreshCookie == nil {
			t.Error("refresh token cookie not found")
		} else {
			if !refreshCookie.HttpOnly {
				t.Errorf("expected HttpOnly to be true, got %v", refreshCookie.HttpOnly)
			}
			if err := validateToken(refreshCookie.Value, "/", "rw"); err != nil {
				t.Errorf("Unexpected error when validating refresh token: %v", err)
			}
		}

	})
}

func TestRefreshHandler(t *testing.T) {
	cfg := config.LoadConfig()
	t.Run("Refresh_Handler_Missing_Cookie", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/refresh", nil)

		RefreshHandler(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected code 401 Unauthorized, received: %d", rr.Code)
		}
		if strings.TrimSpace(rr.Body.String()) != "Unauthorized: No refresh token provided" {
			t.Errorf("Received unexpected error message: %s", strings.TrimSpace(rr.Body.String()))
		}

		respData := TokenResponse{}
		if err := json.Unmarshal(rr.Body.Bytes(), &respData); err == nil {
			t.Errorf("Received access token eventhough refresh token wasn't provided: %s", respData.AccessToken)
		}
	})
	t.Run("Refresh_Handler_Invalid_Token", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/refresh", nil)

		refreshParams := TokenParameters{
			UserId:         "johndoe",
			ExpiryDuration: cfg.Secrets.Jwt.RefreshExpiryDuration,
			FolderId:       "/",
			Access:         "rw",
		}

		refreshClaims := jwt.MapClaims{
			"user_id":   refreshParams.UserId,
			"folder_id": refreshParams.FolderId,
			"access":    refreshParams.Access,
			"exp":       time.Now().Add(refreshParams.ExpiryDuration * time.Hour).Unix(),
		}
		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
		refreshTokenString, err := refreshToken.SignedString([]byte("somewrongjwtinvalidkey"))
		if err != nil {
			t.Fatalf("Received unexpected error when generating token: %v", err)
		}

		req.AddCookie(&http.Cookie{
			Name:     "refresh_token",
			Value:    refreshTokenString,
			Expires:  time.Now().Add(cfg.Secrets.Jwt.RefreshExpiryDuration * time.Hour),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/refresh",
		})

		RefreshHandler(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected code 401 Unauthorized, received: %d", rr.Code)
		}
		if strings.TrimSpace(rr.Body.String()) != "Unauthorized: Invalid refresh token" {
			t.Errorf("Received unexpected error message: %s", strings.TrimSpace(rr.Body.String()))
		}

		respData := TokenResponse{}
		if err := json.Unmarshal(rr.Body.Bytes(), &respData); err == nil {
			t.Errorf("Received access token eventhough refresh token wasn't provided: %s", respData.AccessToken)
		}
	})
	t.Run("Refresh_Handler_Expired_Token", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/refresh", nil)

		refreshParams := TokenParameters{
			UserId:         "johndoe",
			ExpiryDuration: 0,
			FolderId:       "/",
			Access:         "rw",
		}
		_, refreshToken, err := GenerateTokens(&refreshParams, &refreshParams)
		if err != nil {
			t.Fatalf("Received unexpected error when generating token: %v", err)
		}

		req.AddCookie(&http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Expires:  time.Now(),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/refresh",
		})

		RefreshHandler(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected code 401 Unauthorized, received: %d", rr.Code)
		}
		if strings.TrimSpace(rr.Body.String()) != "Unauthorized: Invalid refresh token" {
			t.Errorf("expected `Unauthorized: Invalid refresh token`, got: %s", strings.TrimSpace(rr.Body.String()))
		}

		respData := TokenResponse{}
		if err := json.Unmarshal(rr.Body.Bytes(), &respData); err == nil {
			t.Errorf("Received access token eventhough refresh token wasn't provided: %s", respData.AccessToken)
		}
	})
	t.Run("Refresh_Handler_Success", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/refresh", nil)

		refreshParams := TokenParameters{
			UserId:         "johndoe",
			ExpiryDuration: cfg.Secrets.Jwt.RefreshExpiryDuration,
			FolderId:       "/",
			Access:         "rw",
		}
		_, refreshToken, err := GenerateTokens(&refreshParams, &refreshParams)
		if err != nil {
			t.Fatalf("Received unexpected error when generating token: %v", err)
		}

		req.AddCookie(&http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken, // Use the correct variable
			Expires:  time.Now().Add(cfg.Secrets.Jwt.RefreshExpiryDuration),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/refresh",
		})

		RefreshHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200 OK, got: %d", rr.Code)
		}

		respData := TokenResponse{}
		if err := json.Unmarshal(rr.Body.Bytes(), &respData); err != nil {
			t.Fatalf("error unmarshalling response body: %v", err)
		}
		if respData.AccessToken == "" {
			t.Error("access token not found in response body")
		}

		if err := validateToken(respData.AccessToken, "/", "rw"); err != nil {
			t.Errorf("Unexpected error when validating access token: %v", err)
		}
	})
}

func TestLogoutHandler(t *testing.T) {
	t.Run("Logout_Handler_Success", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/logout", nil)
		rr := httptest.NewRecorder()

		LogoutHandler(rr, req)

		res := rr.Result()
		defer res.Body.Close()

		var refreshCookie *http.Cookie
		for _, cookie := range res.Cookies() {
			if cookie.Name == "refresh_token" {
				refreshCookie = cookie
				break
			}
		}
		if refreshCookie == nil {
			t.Error("expected refresh_token cookie to be set")
		} else {
			if refreshCookie.Value != "" {
				t.Errorf("expected cookie value to be empty, got '%s'", refreshCookie.Value)
			}
			if !refreshCookie.Expires.Equal(time.Unix(0, 0)) {
				t.Errorf("expected cookie Expires to be %v, got %v", time.Unix(0, 0), refreshCookie.Expires)
			}
			if refreshCookie.HttpOnly != true {
				t.Errorf("expected HttpOnly to be true, got %v", refreshCookie.HttpOnly)
			}
		}

		var tokenResponse TokenResponse
		if err := json.NewDecoder(res.Body).Decode(&tokenResponse); err != nil {
			t.Errorf("error decoding response body: %v", err)
		}
		if tokenResponse.AccessToken != "" {
			t.Errorf("expected AccessToken to be empty, got '%s'", tokenResponse.AccessToken)
		}
	})
}

func TestSharingGatewayWrongArguments(t *testing.T) {
	
	url := "/auth-share"
	expirationDate := time.Now().Add(48 * time.Hour).UTC().Truncate(time.Second)
	expiration := expirationDate.Format(time.RFC3339)
	otpPass := "123456"
	linkUrl := uuid.New().String()
	expiryDuration := expirationDate.Sub(time.Now().UTC().Truncate(time.Second))
	sharingFolderId := helpers.GenerateFolderName(expiryDuration, linkUrl)
	folderName := "someFolderName"
	access := "rw"
	salt, err := helpers.GenerateRandomSalt()
	if err != nil {
		t.Errorf("Received unexpected error when generating random salt: %v", err)
	}

	db, mock, err := initMockDb()
	if err != nil {
		t.Fatalf("Received unexpected error when initializing mock db: %v", err)
	}
	defer db.Close()

	t.Run("Non_Existent_LinkUrl", func(t *testing.T) {
		wrongLinkUrl := uuid.New().String()
		expectedResponse :=  "Forbidden: user not found" 
		rows := sqlmock.NewRows([]string{"link_url", "folder_id", "folder_name", "salt", "otp_hash", "access", "expiration"})
	
		mock.ExpectQuery("SELECT link_url, folder_id, folder_name, salt, otp_hash, access, expiration FROM sharing_users WHERE link_url = \\$1").
			WithArgs(wrongLinkUrl).
			WillReturnRows(rows)

		sharingCreds := SharingCredentials{
			LinkUrl: wrongLinkUrl,
			OtpPassword: otpPass,
		}
		body, err := json.Marshal(sharingCreds)
		if err != nil {
			t.Fatalf("failed to marshal credentials: %v", err)
		}
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("POST", url, bytes.NewBuffer(body))

		SharingGatewayHandler(rr, req, db)

		if rr.Code != http.StatusForbidden {
			t.Errorf("Expected status 403 Forbidden, got: %d", rr.Code)
		}

		if strings.TrimSpace(rr.Body.String()) != expectedResponse {
			t.Errorf("Expected response %s, got: %s", expectedResponse, strings.TrimSpace(rr.Body.String()))
		}
	})

	t.Run("Wrong_Otp_Pass", func(t *testing.T) {
		expectedResponse :=  "Forbidden: invalid credentials"
		wrongPass := "000000"
		wrongHash := helpers.HashPassword(wrongPass, salt)
		rows := sqlmock.NewRows([]string{"link_url", "folder_id", "folder_name", "salt", "otp_hash", "access", "expiration"}).
					AddRow(linkUrl, sharingFolderId, folderName, salt, wrongHash, access, expiration)
	
		mock.ExpectQuery("SELECT link_url, folder_id, folder_name, salt, otp_hash, access, expiration FROM sharing_users WHERE link_url = \\$1").
			WithArgs(linkUrl).
			WillReturnRows(rows)

		sharingCreds := SharingCredentials{
			LinkUrl: linkUrl,
			OtpPassword: otpPass,
		}
		body, err := json.Marshal(sharingCreds)
		if err != nil {
			t.Fatalf("failed to marshal credentials: %v", err)
		}
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("POST", url, bytes.NewBuffer(body))

		SharingGatewayHandler(rr, req, db)

		if rr.Code != http.StatusForbidden {
			t.Errorf("Expected status 403 Forbidden, got: %d", rr.Code)
		}

		if strings.TrimSpace(rr.Body.String()) != expectedResponse {
			t.Errorf("Expected response %s, got: %s", expectedResponse, strings.TrimSpace(rr.Body.String()))
		}
	})

	
}

func TestSharingGatewayHandlerSucces(t *testing.T) {

	url := "/auth-share"

	expirationDate := time.Now().Add(48 * time.Hour).UTC().Truncate(time.Second)
	expiration := expirationDate.Format(time.RFC3339)
	otpPass := "123456"
	linkUrl := uuid.New().String()
	expiryDuration := expirationDate.Sub(time.Now().UTC().Truncate(time.Second))
	sharingFolderId := helpers.GenerateFolderName(expiryDuration, linkUrl)
	folderName := "someFolderName"
	access := "w"
	salt, err := helpers.GenerateRandomSalt()
	if err != nil {
		t.Errorf("Received unexpected error when generating random salt: %v", err)
	}
	hashedOtp := helpers.HashPassword(otpPass, salt)

	db, mock, err := initMockDb()
	if err != nil {
		t.Fatalf("Received unexpected error when initializing mock db: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"link_url", "folder_id", "folder_name", "salt", "otp_hash", "access", "expiration"}).
			AddRow(linkUrl, sharingFolderId, folderName, salt, hashedOtp, access, expiration)

	mock.ExpectQuery("SELECT link_url, folder_id, folder_name, salt, otp_hash, access, expiration FROM sharing_users WHERE link_url = \\$1").
		WithArgs(linkUrl).
		WillReturnRows(rows)

	sharingCreds := SharingCredentials{
		LinkUrl: linkUrl,
		OtpPassword: otpPass,
	}
	body, err := json.Marshal(sharingCreds)
	if err != nil {
		t.Fatalf("failed to marshal credentials: %v", err)
	}
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", url, bytes.NewBuffer(body))

	SharingGatewayHandler(rr, req, db)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got: %d", rr.Code)
	}

	// Check if access token, refresh token is set and valid.
	respData := SharingTokenResponse{}
	if err := json.Unmarshal(rr.Body.Bytes(), &respData); err != nil {
		t.Fatalf("error unmarshalling response body: %v", err)
	}
	if respData.AccessToken == "" {
		t.Error("access token not found in response body")
	}
	if respData.FolderId != sharingFolderId {
		t.Errorf("Expected folderId : %s, got : %s", sharingFolderId, respData.FolderId)
	}

	if err := validateToken(respData.AccessToken, sharingFolderId, access); err != nil {
		t.Errorf("Unexpected error when validating access token: %v", err)
	}

	var refreshCookie *http.Cookie
	for _, cookie := range rr.Result().Cookies() {
		if cookie.Name == "refresh_token" {
			refreshCookie = cookie
			break
		}
	}
	if refreshCookie == nil {
		t.Error("refresh token cookie not found")
	} else {
		if !refreshCookie.HttpOnly {
			t.Errorf("expected HttpOnly to be true, got %v", refreshCookie.HttpOnly)
		}
		if err := validateToken(refreshCookie.Value, sharingFolderId, access); err != nil {
			t.Errorf("Unexpected error when validating refresh token: %v", err)
		}
	}
}