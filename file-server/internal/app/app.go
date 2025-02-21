package app

import (
	"file-server/internal/uploader"
	"net/http"

	"github.com/rs/cors"
)

func SetupServer() *http.Server {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://kuza.gr"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/upload", uploader.UploadHandler)

	return &http.Server{
		Addr:    ":8080",
		Handler: c.Handler(mux),
	}
}