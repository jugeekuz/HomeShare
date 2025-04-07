// main.go
package main

import (
	"log"
	"os"
	"time"

	"file-server/internal/app"
	"file-server/config"
	"file-server/internal/helpers"
	"file-server/internal/job"
)

func main() {
	cfg := config.LoadConfig()
	if err := os.MkdirAll(cfg.SharingDir, os.ModePerm); err != nil {
		log.Fatal("[FILE-SERVER] Error while creating sharing folder directory")
	}

	go func () {
		_ = helpers.CleanupExpiredFolders(cfg.SharingDir)
		time.Sleep(30 * time.Minute)
	}()

	job_timeout := 45 * time.Second
	jm := job.NewJobManager(job_timeout)

	server, err := app.SetupServer(jm, app.InitDatabase)
	if err != nil {
		log.Fatalf("[FILE-SERVER] Server setup failed: %v", err)
	}

	server.Addr = ":443"

	log.Fatal(server.ListenAndServeTLS(
		"./certs/fullchain.pem",
		"./certs/privkey.pem",
	))
}
