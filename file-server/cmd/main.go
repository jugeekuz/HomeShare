// main.go
package main

import (
	"log"
	"time"
	"file-server/internal/app"
	"file-server/internal/job"
)

func main() {

	job_timeout := 45 * time.Second
	jm := job.NewJobManager(job_timeout)

	server := app.SetupServer(jm)

	server.Addr = ":443"

	log.Fatal(server.ListenAndServeTLS(
		"./certs/fullchain.pem",
		"./certs/privkey.pem",
	))
}
