// main.go
package main

import (
	"file-server/internal/app"
	"file-server/config"
	"file-server/internal/helpers"
	"file-server/internal/job"
	"log"
	"os"
	"time"
)

func main() {
	cfg := config.LoadConfig()
	if err := os.MkdirAll(cfg.SharingDir, os.ModePerm); err != nil {
		log.Fatal("Error while creating sharing folder directory")
	}

	go func () {
		_ = helpers.CleanupExpiredFolders(cfg.SharingDir)
		time.Sleep(30 * time.Minute)
	}()

	job_timeout := 45 * time.Second
	jm := job.NewJobManager(job_timeout)

	server := app.SetupServer(jm)

	server.Addr = ":443"

	log.Fatal(server.ListenAndServeTLS(
		"./certs/fullchain.pem",
		"./certs/privkey.pem",
	))
}
