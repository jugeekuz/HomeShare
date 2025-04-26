package repositories

import (
	"fmt"
	"errors"
	"database/sql"

	"file-server/internal/helpers"
	"file-server/internal/models"
)

func InitializeSharingUserTable(db *sql.DB) error {
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS sharing_users (
			link_url TEXT PRIMARY KEY,
			folder_id TEXT NOT NULL,
			folder_name TEXT NOT NULL,
			salt TEXT NOT NULL,
			otp_hash TEXT NOT NULL,
			access TEXT NOT NULL CHECK (access IN ('r', 'w', 'rw')),
			expiration TIMESTAMPTZ NOT NULL
		)
	`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("error creating sharing_users table: %w", err)
	}
	createExpIndexQuery := `
		CREATE INDEX IF NOT EXISTS idx_expiration ON sharing_users (expiration);
	`
	_, err = db.Exec(createExpIndexQuery)
	if err != nil {
		return fmt.Errorf("error creating sharing_users expiration index: %w", err)
	}

	return nil
}

func CreateSharingUser(db *sql.DB, linkUrl string, folderId string, folderName string, salt string, otpPass string, access string, expiration string) (models.SharingUser, error) {
	var sharingUser models.SharingUser

	createUserQuery := `
		INSERT INTO sharing_users (link_url, folder_id, folder_name, salt, otp_hash, access, expiration)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (link_url) DO UPDATE 
		SET link_url = EXCLUDED.link_url 
		RETURNING link_url, folder_id, folder_name, salt, otp_hash, access, expiration;
	`

	sharingUser.LinkUrl = linkUrl
	sharingUser.FolderId = folderId
	sharingUser.FolderName = folderName
	sharingUser.Salt = salt
	sharingUser.OtpHash = helpers.HashPassword(otpPass, salt)
	sharingUser.Access = access
	sharingUser.Expiration = expiration

	_, err := db.Exec(createUserQuery, sharingUser.LinkUrl, sharingUser.FolderId, sharingUser.FolderName, sharingUser.Salt, sharingUser.OtpHash, sharingUser.Access, sharingUser.Expiration)
	if err != nil {
		return models.SharingUser{}, err
	}
	return sharingUser, nil
}

func GetSharingUser(db *sql.DB, linkUrl string) (*models.SharingUser, error) {
	query := `
		SELECT link_url, folder_id, folder_name, salt, otp_hash, access, expiration
		FROM sharing_users
		WHERE link_url = $1
	`
	row := db.QueryRow(query, linkUrl)
	var user models.SharingUser
	err := row.Scan(&user.LinkUrl, &user.FolderId, &user.FolderName, &user.Salt, &user.OtpHash, &user.Access, &user.Expiration)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func DeleteSharingUser(db *sql.DB, linkUrl string) error {
	query := `DELETE FROM sharing_users WHERE link_url = $1`
	result, err := db.Exec(query, linkUrl)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func DeleteExpiredUsers(db *sql.DB) error {
	query := "DELETE FROM sharing_users WHERE expiration < NOW()"

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	// go func() {
	// 	rowsAffected, err := res.RowsAffected()
	// 	if err != nil {
	// 		fmt.Printf("Error retrieving affected rows: %v", err)
	// 	}
	// 	fmt.Printf("Deleted %d rows.\n", rowsAffected)
	// }()
	return nil
}
