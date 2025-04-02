package auth

import (
	"fmt"
	"database/sql"
	"encoding/json"
	"file-server/config"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenParameters struct {
	UserId         string        `json:"user_id"`
	ExpiryDuration time.Duration `json:"exp"`
	FolderId       string        `json:"folder_id"`
	Access         string        `json:"access"` // "r", "w", or "rw"
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	cfg := config.LoadConfig()

	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if _, err := Authenticate(db, creds); err != nil {
		http.Error(w, fmt.Sprintf("Forbidden: %v", err), http.StatusForbidden)
		return
	}

	accessParams := &TokenParameters{
		UserId:         creds.Username,
		ExpiryDuration: cfg.Secrets.Jwt.AccessExpiryDuration,
		FolderId:       "/",
		Access:         "rw",
	}
	refreshParams := &TokenParameters{
		UserId:         creds.Username,
		ExpiryDuration: cfg.Secrets.Jwt.RefreshExpiryDuration,
		FolderId:       "/",
		Access:         "rw",
	}

	accessTokenString, refreshTokenString, err := GenerateTokens(accessParams, refreshParams)
	if err != nil {
		http.Error(w, "Issue generating tokens", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshTokenString,
		Expires:  time.Now().Add(cfg.Secrets.Jwt.RefreshExpiryDuration),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	})

	response := TokenResponse{
		AccessToken: accessTokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	cfg := config.LoadConfig()
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Unauthorized: No refresh token provided", http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secrets.Jwt.JwtSecret), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized: Invalid refresh token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized: Invalid token claims", http.StatusUnauthorized)
		return
	}

	userId, ok1 := claims["user_id"].(string)
	folderId, ok2 := claims["folder_id"].(string)
	access, ok3 := claims["access"].(string)
	if !ok1 || !ok2 || !ok3 {
		http.Error(w, "Unauthorized: Malformed token", http.StatusUnauthorized)
		return
	}

	accessParams := &TokenParameters{
		UserId:         userId,
		ExpiryDuration: cfg.Secrets.Jwt.AccessExpiryDuration,
		FolderId:       folderId,
		Access:         access,
	}

	accessTokenString, _, err := GenerateTokens(accessParams, accessParams)
	if err != nil {
		http.Error(w, "Could not generate new access token", http.StatusInternalServerError)
		return
	}

	response := TokenResponse{
		AccessToken: accessTokenString,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0), // Expire immediately
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	})

	response := TokenResponse{
		AccessToken: "",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
