package app

import (
	// "context"
	"database/sql"
	"file-server/config"
	"file-server/internal/auth"
	"file-server/internal/db"
	"file-server/internal/downloader"
	"file-server/internal/job"
	"file-server/internal/sharing"
	"file-server/internal/uploader"
	"file-server/internal/repositories"
	"file-server/internal/helpers"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/cors"
	"github.com/google/uuid"
)

type DatabaseCallback func() (*sql.DB, error)

func InitDatabase() (*sql.DB, error) {
	cfg := config.LoadConfig()

	db, err := db.StartConn()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	user := cfg.User
	if err := repositories.InitializeUserTable(db); err != nil {
		return nil, err
	}
	if err := repositories.InitializeSharingUserTable(db); err != nil {
		return nil, err
	}
	if _, err := repositories.CreateAdminUser(db, user.Username, user.Email, user.Password); err != nil {
		return nil, err
	}
	return db, nil
}

func SetupServer(jm *job.JobManager, dbCallback DatabaseCallback) (*http.Server, error) {
	cfg := config.LoadConfig()

	db, err := dbCallback()
	if err != nil {
		return nil, err
	}
	
	go func(){
		for {
			err := repositories.DeleteExpiredUsers(db)
			if err != nil {
				fmt.Printf("Received unexpected error when deleting expired users")
			}
			time.Sleep(5*time.Second)
		}
	}()
	

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.DomainOrigin, "http://localhost:3001"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Set-Cookie", "Folder-Id"},
		AllowCredentials: true,
	})

	mux := http.NewServeMux()

	// Unauthenticated endpoints
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		auth.LoginHandler(w, r, db)
	})

	mux.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		auth.RefreshHandler(w, r)
	})

	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		auth.LogoutHandler(w, r)
	})

	mux.HandleFunc("/auth-share", func(w http.ResponseWriter, r *http.Request) {
		auth.SharingGatewayHandler(w, r, db)
	})

	// Authenticated endpoints
	mux.HandleFunc("/upload",
		auth.AuthMiddleware(
			func(w http.ResponseWriter, r *http.Request) {
				uploader.UploadHandler(w, r, jm, cfg.UploadDir)
			}))

	mux.HandleFunc("/download",
		auth.RefreshAuthMiddleware(
			func(w http.ResponseWriter, r *http.Request) {
				downloader.DownloadHandler(w, r, jm)
			}))

	mux.HandleFunc("/download-available",
		auth.AuthMiddleware(
			func(w http.ResponseWriter, r *http.Request) {
				downloader.DownloadAvailableHandler(w, r, jm)
			}))

	mux.HandleFunc("/share",
		auth.AuthMiddleware(
			func(w http.ResponseWriter, r *http.Request) {
				salt, err := helpers.GenerateRandomSalt()
				if err != nil {
					http.Error(w, fmt.Sprintf("Error while creating user: %v", err), http.StatusInternalServerError)
					return
				}
				sharing.SharingHandler(w, r, db, salt, uuid.New().String())
			}))

	mux.HandleFunc("/share-file",
		auth.AuthMiddleware(
			func(w http.ResponseWriter, r *http.Request) {
				sharing.AddSharingFilesHandler(w, r, jm)
			}))

	mux.HandleFunc("/share-files",
		auth.AuthMiddleware(
			func(w http.ResponseWriter, r *http.Request) {
				sharing.GetSharingFilesHandler(w, r)
			}))

	return &http.Server{
		Addr:    ":443",
		Handler: c.Handler(mux),
	}, nil
}
