package config

import (
	"os"
	"log"
	"strconv"
	"time"
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

type JWT struct {
	JwtSecret				string
	AccessExpiryDuration	time.Duration	
	RefreshExpiryDuration	time.Duration	
}

type Secrets struct {
	Jwt		JWT
}
type Config struct {
	DomainOrigin	string
	UploadDir       string
	SharingDir      string
	ChunksDir 	    string
	DB 				DBConfig
	Secrets			Secrets
	User			User
}

func LoadConfig() *Config {
	accessTokenExp, err := strconv.Atoi(getEnv("ACCESS_TOKEN_EXP_H", "1"))
	if err != nil {
		log.Fatalf("Invalid ACCESS_TOKEN_EXP_H value: %v", err)
	}
	refreshTokenExp, err := strconv.Atoi(getEnv("REFRESH_TOKEN_EXP_H", "1"))
	if err != nil {
		log.Fatalf("Invalid REFRESH_TOKEN_EXP_H value: %v", err)
	}
	return &Config{
		DomainOrigin:	getEnv("DOMAIN_ORIGIN", "https://kuza.gr"),
		UploadDir: 		getEnv("UPLOAD_DIR", "uploads"),
		SharingDir: 	getEnv("SHARING_DIR", "temp"),
		ChunksDir: 		getEnv("CHUNKS_DIR", "chunks"),
		DB: DBConfig{
            Host:     	getEnv("DB_HOST", "postgres"),
            Port:     	getEnv("DB_PORT", "5432"),
            User:     	getEnv("POSTGRES_USER", "myuser"),
            Password: 	getEnv("POSTGRES_PASSWORD", "mypassword"),
            DBName:   	getEnv("POSTGRES_DB", "userdb"),
            SSLMode:  	getEnv("DB_SSL_MODE", "disable"),
        },
		Secrets: Secrets{
			JWT {
				JwtSecret: 				GetOrCreateJWTSecret("secrets", "JWT"),
				AccessExpiryDuration: 	time.Duration(accessTokenExp),
				RefreshExpiryDuration: 	time.Duration(refreshTokenExp),
			},
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
