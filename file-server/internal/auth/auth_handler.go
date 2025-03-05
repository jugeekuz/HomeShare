package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenParameters struct {
	UserId         string        `json:"user_id"`
	ExpiryDuration time.Duration `json:"expiry_duration_hours"`
	FolderId       string        `json:"folder_id"`
	Access         string        `json:"access"` // "read", "write", or "read+write"
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func GenerateTokens(accessParams, refreshParams *TokenParameters) (string, string, error) {
	accessClaims := jwt.MapClaims{
		"user_id":   accessParams.UserId,
		"folder_id": accessParams.FolderId,
		"access":    accessParams.Access,
		"exp":       time.Now().Add(accessParams.ExpiryDuration * time.Hour).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	refreshClaims := jwt.MapClaims{
		"user_id":   refreshParams.UserId,
		"folder_id": refreshParams.FolderId,
		"access":    refreshParams.Access,
		"exp":       time.Now().Add(refreshParams.ExpiryDuration * time.Hour).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	accessParams := &TokenParameters{
		UserId:         creds.Username,
		ExpiryDuration: 1, // Access token expires in 1 hour.
		FolderId:       "/",
		Access:         "rw",
	}
	refreshParams := &TokenParameters{
		UserId:         creds.Username,
		ExpiryDuration: 24, // Refresh token expires in 24 hours.
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
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/refresh",
	})

	response := TokenResponse{
		AccessToken: accessTokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Unauthorized: No refresh token provided", http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
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
		ExpiryDuration: 1,
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
		SameSite: http.SameSiteStrictMode,
		Path:     "/refresh",
	})

	response := TokenResponse{
		AccessToken: "",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
