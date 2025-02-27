package app

import (
	"file-server/internal/job"
	"file-server/internal/uploader"
	"net/http"

	"github.com/rs/cors"
)

func SetupServer(jm *job.JobManager) *http.Server {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://kuza.gr", "http://kuza.gr/", "https://kuza.gr", "https://kuza.gr/"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		uploader.UploadHandler(w, r, jm)
	})

	return &http.Server{
		Addr:    ":8080",
		Handler: c.Handler(mux),
	}
}
