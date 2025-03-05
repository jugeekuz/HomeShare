package app

import (
	"file-server/internal/job"
	"file-server/internal/uploader"
	"file-server/internal/downloader"
	"file-server/internal/auth"
	"net/http"

	"github.com/rs/cors"
)

func SetupServer(jm *job.JobManager) *http.Server {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"https://kuza.gr", "https://kuza.gr/"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		uploader.UploadHandler(w, r, jm)
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

	return &http.Server{
		Addr:    ":443",
		Handler: c.Handler(mux),
	}
}
