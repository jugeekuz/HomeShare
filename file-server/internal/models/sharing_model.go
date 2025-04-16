package models

type SharingUser struct {
	LinkUrl			string `json:"link_url"`
	FolderId		string `json:"folder_id"`
	FolderName		string `json:"folder_name"`
	Salt			string `json:"salt"`
	OtpHash			string `json:"otp_hash"`
	Access			string `json:"access"`
	Expiration		string `json:"expiration"`
}