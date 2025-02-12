// main.go
package main

import (
	"log"
	"file-server/internal/app"
)

func main() {
	server := app.SetupServer()
	log.Fatal(server.ListenAndServe())
}
