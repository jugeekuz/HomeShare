package models

type User struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	Salt         string `json:"salt"`
	PasswordHash string `json:"password_hash"`
	FolderId     string `json:"folder_id"`
	Access       string `json:"access"` // 'r', 'w' or 'rw'
}

