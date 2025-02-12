package app

import (
    "net/http"
    "github.com/rs/cors"
	"file-server/internal/uploader"
)

func SetupServer() *http.Server {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://kuza.gr"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
	})
	
	mux := http.NewServeMux()
	mux.HandleFunc("/uploads", uploader.UploadHandler)

	return &http.Server{
        Addr:    ":8080",
        Handler: c.Handler(mux),
    }
}