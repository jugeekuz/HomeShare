package auth_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"file-server/internal/auth"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("my_secret_key")

func TestGenerateTokens(t *testing.T) {
	accessParams := &auth.TokenParameters{
		UserId:         "testuser",
		ExpiryDuration: 1, // 1 hour
		FolderId:       "/",
		Access:         "rw",
	}
	refreshParams := &auth.TokenParameters{
		UserId:         "testuser",
		ExpiryDuration: 24, // 24 hours
		FolderId:       "/",
		Access:         "rw",
	}

	accessTokenStr, refreshTokenStr, err := auth.GenerateTokens(accessParams, refreshParams)
	if err != nil {
		t.Fatalf("GenerateTokens returned error: %v", err)
	}
	if accessTokenStr == "" || refreshTokenStr == "" {
		t.Fatal("Expected non-empty token strings")
	}

	// Parse access token to check claims.
	token, err := jwt.Parse(accessTokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		t.Fatalf("Invalid access token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Expected token claims to be of type jwt.MapClaims")
	}
	if claims["user_id"] != "testuser" {
		t.Errorf("Expected user_id to be 'testuser', got %v", claims["user_id"])
	}
	if claims["folder_id"] != "/" {
		t.Errorf("Expected folder_id to be '/', got %v", claims["folder_id"])
	}
	if claims["access"] != "rw" {
		t.Errorf("Expected access to be 'rw', got %v", claims["access"])
	}
}

func TestLoginHandler(t *testing.T) {
	creds := auth.Credentials{
		Username: "testuser",
		Password: "password",
	}
	bodyBytes, _ := json.Marshal(creds)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(bodyBytes))
	rec := httptest.NewRecorder()

	auth.LoginHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", res.StatusCode)
	}

	cookies := res.Cookies()
	var found bool
	for _, cookie := range cookies {
		if cookie.Name == "refresh_token" {
			found = true
			if cookie.Value == "" {
				t.Error("Expected non-empty refresh_token cookie")
			}
		}
	}
	if !found {
		t.Fatal("Expected refresh_token cookie to be set")
	}

	var tokenResp auth.TokenResponse
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		t.Fatalf("Error unmarshaling response JSON: %v", err)
	}
	if tokenResp.AccessToken == "" {
		t.Error("Expected non-empty access token in response")
	}
}

func TestLoginHandlerInvalidPayload(t *testing.T) {
	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString("invalid json"))
	rec := httptest.NewRecorder()

	auth.LoginHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status BadRequest, got %v", res.StatusCode)
	}
}
func TestRefreshHandler(t *testing.T) {
	refreshParams := &auth.TokenParameters{
		UserId:         "testuser",
		ExpiryDuration: 24,
		FolderId:       "/",
		Access:         "rw",
	}
	_, refreshTokenStr, err := auth.GenerateTokens(refreshParams, refreshParams)
	if err != nil {
		t.Fatalf("Error generating tokens: %v", err)
	}

	req := httptest.NewRequest("POST", "/refresh", nil)
	req.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: refreshTokenStr,
	})
	rec := httptest.NewRecorder()

	auth.RefreshHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", res.StatusCode)
	}

	var tokenResp auth.TokenResponse
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		t.Fatalf("Error unmarshaling response JSON: %v", err)
	}
	if tokenResp.AccessToken == "" {
		t.Error("Expected non-empty access token in response")
	}
}
func TestRefreshHandlerInvalidCookie(t *testing.T) {
	req := httptest.NewRequest("POST", "/refresh", nil)
	rec := httptest.NewRecorder()

	auth.RefreshHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected status Unauthorized, got %v", res.StatusCode)
	}
}
func TestLogoutHandler(t *testing.T) {
	req := httptest.NewRequest("POST", "/logout", nil)
	rec := httptest.NewRecorder()

	auth.LogoutHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	cookies := res.Cookies()
	var found bool
	for _, cookie := range cookies {
		if cookie.Name == "refresh_token" {
			found = true
			if cookie.Value != "" {
				t.Error("Expected refresh_token cookie value to be empty on logout")
			}
			if cookie.Expires.After(time.Now()) {
				t.Error("Expected refresh_token cookie to be expired")
			}
		}
	}
	if !found {
		t.Fatal("Expected refresh_token cookie to be set on logout")
	}
	var tokenResp auth.TokenResponse
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		t.Fatalf("Error unmarshaling response JSON: %v", err)
	}
	if tokenResp.AccessToken != "" {
		t.Error("Expected access token to be empty in logout response")
	}
}
