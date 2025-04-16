package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	
	"file-server/internal/helpers"
	"file-server/internal/models"

)

type UserRepository interface {
}

type userRepository struct {
	db *sql.DB
}

func InitializeUserTable(db *sql.DB) error {
	createTableQuery :=
		`
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

func CreateAdminUser(db *sql.DB, username string, email string, password string) (models.User, error) {
	var user models.User

	createUserQuery :=
		`
		INSERT INTO users (username, email, salt, password_hash, folder, access)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (username) DO UPDATE 
		SET email = EXCLUDED.email 
		RETURNING username, email, salt, password_hash, folder, access;
	`
	salt, err := helpers.GenerateRandomSalt()
	if err != nil {
		return models.User{}, err
	}
	user.Username = username
	user.Email = email
	user.Salt = salt
	user.PasswordHash = helpers.HashPassword(password, salt)
	user.FolderId = "/"
	user.Access = "rw"

	_, err = db.Exec(createUserQuery, user.Username, user.Email, user.Salt, user.PasswordHash, user.FolderId, user.Access)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func GetUserByUsername(db *sql.DB, username string) (*models.User, error) {
	query :=
		`
		SELECT username, email, salt, password_hash, folder, access
		FROM users
		WHERE username = $1
	`
	row := db.QueryRow(query, username)
	var user models.User
	err := row.Scan(&user.Username, &user.Email, &user.Salt, &user.PasswordHash, &user.FolderId, &user.Access)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
