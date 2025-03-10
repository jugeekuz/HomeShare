package app

import (
	"file-server/config"
	"file-server/internal/job"
	"file-server/internal/uploader"
	"file-server/internal/sharing"
	"file-server/internal/downloader"
	"file-server/internal/auth"
	"net/http"

	"github.com/rs/cors"
)

func SetupServer(jm *job.JobManager) *http.Server {
	cfg := config.LoadConfig()
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"https://kuza.gr", "https://kuza.gr/"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		uploader.UploadHandler(w, r, cfg.UploadDir)
	})

	mux.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		downloader.DownloadHandler(w, r, jm)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		auth.LoginHandler(w, r)
	})

	mux.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		auth.RefreshHandler(w, r)
	})

	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		auth.LogoutHandler(w, r)
	})

	mux.HandleFunc("/share", func(w http.ResponseWriter, r *http.Request) {
		sharing.SharingHandler(w, r)
	})

	mux.HandleFunc("/share-file", func(w http.ResponseWriter, r *http.Request) {
		sharing.AddSharingFilesHandler(w, r, jm)
	})

	return &http.Server{
		Addr:    ":443",
		Handler: c.Handler(mux),
	}
}
