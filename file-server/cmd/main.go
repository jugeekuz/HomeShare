// main.go
package main

import (
	// "bufio"
	"fmt"
	"log"
	"net/http"
	"file-server/internal/handlers"
)


func main() {
	http.HandleFunc("/upload", uploader.UploadHandler)
	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
