package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"file-server/config"

	"github.com/golang-jwt/jwt/v5"
)
type contextKey string

const ClaimsContextKey contextKey = "claims"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := config.LoadConfig()
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}
		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.Secrets.Jwt.JwtSecret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsContextKey, token.Claims)
		next(w, r.WithContext(ctx))
	})
}

func RefreshAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		ctx := context.WithValue(r.Context(), ClaimsContextKey, token.Claims)
		next(w, r.WithContext(ctx))
	})
}

func CookieHasAccess(r *http.Request, folderId string, access string) bool {
	cfg := config.LoadConfig()
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		return false
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secrets.Jwt.JwtSecret), nil
	})

	if err != nil || !token.Valid {
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	_, ok1 := claims["user_id"].(string)
	folderId, ok2 := claims["folder_id"].(string)
	access, ok3 := claims["access"].(string)
	if !ok1 || !ok2 || !ok3 {
		return false
	}

	hasAccess, err := HasAccess(claims, folderId, access)
	if !hasAccess || err != nil {
		return false
	}

	return true
}
