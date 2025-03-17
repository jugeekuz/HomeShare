package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"file-server/config"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Username     string
	Email        string
	Salt         string
	PasswordHash string
	FolderId     string
	Access       string // r, w or rw
}

type contextKey string

const ClaimsContextKey contextKey = "claims"

func HashPassword(password, salt string) string {
	hash := sha256.Sum256([]byte(password + salt))
	return hex.EncodeToString(hash[:])
}

func GenerateRandomSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

func InitializeUserTable(db *sql.DB) error {
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS users (
			username TEXT PRIMARY KEY,
			email TEXT NOT NULL,
			salt TEXT NOT NULL,
			password_hash TEXT NOT NULL,
			folder TEXT NOT NULL,
			access TEXT NOT NULL CHECK (access IN ('r', 'w', 'rw'))
		)
	`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}
	return nil
}

func CreateAdminUser(db *sql.DB, username string, email string, password string) (User, error) {
	var user User

	createUserQuery := `
		INSERT INTO users (username, email, salt, password_hash, folder, access)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (username) DO UPDATE 
		SET email = EXCLUDED.email 
		RETURNING username, email, salt, password_hash, folder, access;
	`
	salt, err := GenerateRandomSalt()
	if err != nil {
		return User{}, err
	}
	user.Username = username
	user.Email = email
	user.Salt = salt
	user.PasswordHash = HashPassword(password, salt)
	user.FolderId = "/"
	user.Access = "rw"

	_, err = db.Exec(createUserQuery, user.Username, user.Email, user.Salt, user.PasswordHash, user.FolderId, user.Access)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func getUserByUsername(db *sql.DB, username string) (*User, error) {
	query := `
		SELECT username, email, salt, password_hash, folder, access
		FROM users
		WHERE username = $1
	`
	row := db.QueryRow(query, username)
	var user User
	err := row.Scan(&user.Username, &user.Email, &user.Salt, &user.PasswordHash, &user.FolderId, &user.Access)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

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
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsContextKey, token.Claims)
		next(w, r.WithContext(ctx))
	})
}
