package auth

import (
	"fmt"
	"errors"
	"crypto/rand"
	"database/sql"
	"crypto/sha256"
	"encoding/hex"
)
type User struct {
	Username		string
	Email			string
	Salt			string
	PasswordHash	string
	FolderId		string
	Access			string // r, w or rw
}

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

