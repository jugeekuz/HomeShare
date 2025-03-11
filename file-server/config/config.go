package config

import (
	"os"

	"file-server/internal/helpers"
)
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type User struct {
	Username	string
	Email		string
	Password	string
}

type Secrets struct {
	JwtSecret	string	
}
type Config struct {
	UploadDir       string
	SharingDir      string
	DB 				DBConfig
	Secrets			Secrets
	User			User
}

func LoadConfig() *Config {
	return &Config{
		UploadDir: 		getEnv("UPLOAD_DIR", "uploads"),
		SharingDir: 	getEnv("SHARING_DIR", "temp"),
		DB: DBConfig{
            Host:     	getEnv("DB_HOST", "postgres"),
            Port:     	getEnv("DB_PORT", "5432"),
            User:     	getEnv("POSTGRES_USER", "myuser"),
            Password: 	getEnv("POSTGRES_PASSWORD", "mypassword"),
            DBName:   	getEnv("POSTGRES_DB", "userdb"),
            SSLMode:  	getEnv("DB_SSL_MODE", "disable"),
        },
		Secrets: Secrets{
			JwtSecret: helpers.GetOrCreateJWTSecret("secrets", "JWT"),
		},
		User: User{
			Username: 	getEnv("ADMIN_USERNAME", "admin"),
			Password: 	getEnv("ADMIN_PASSWORD", "admin"),
			Email:	  	getEnv("ADMIN_EMAIL", "admin@email.com"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
