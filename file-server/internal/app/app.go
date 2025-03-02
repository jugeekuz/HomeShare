package app

import (
	"file-server/internal/job"
	"file-server/internal/uploader"
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

	return &http.Server{
		Addr:    ":443",
		Handler: c.Handler(mux),
	}
}
