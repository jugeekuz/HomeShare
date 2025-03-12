package app

import (
	// "context"
	"file-server/config"
	"file-server/internal/auth"
	"file-server/internal/db"
	"file-server/internal/downloader"
	"file-server/internal/job"
	"file-server/internal/sharing"
	"file-server/internal/uploader"
	"fmt"
	"net/http"

	"github.com/rs/cors"
)

func SetupServer(jm *job.JobManager) (*http.Server, error) {
	cfg := config.LoadConfig()

	db, err := db.StartConn()
	if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }
	user := cfg.User
	if err := auth.InitializeUserTable(db); err != nil {
		return nil, err
	}
	if _, err := auth.CreateAdminUser(db, user.Username, user.Email, user.Password); err != nil {
		return nil, err
	}

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"https://kuza.gr", "https://kuza.gr/"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
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

	// Authenticated endpoints
	mux.HandleFunc("/upload", 
		auth.AuthMiddleware(
			func(w http.ResponseWriter, r *http.Request) {
				uploader.UploadHandler(w, r, cfg.UploadDir)
			}))

	mux.HandleFunc("/download", 
		auth.AuthMiddleware(
			func(w http.ResponseWriter, r *http.Request) {
				downloader.DownloadHandler(w, r, jm)
			}))

	mux.HandleFunc("/share", 
		auth.AuthMiddleware(
			func(w http.ResponseWriter, r *http.Request) {
				sharing.SharingHandler(w, r)
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
