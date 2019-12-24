package main

import (
	handler "github.com/sil-vio/golang-file-server/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler.UploadHandler)
	http.HandleFunc("/upload", handler.UploadHandler)
	http.HandleFunc("/download", handler.DownloadHandler)
	log.Print("Listening on port:4200...")
	// Listen on port 8080
	err := http.ListenAndServe(":4200", nil)
	if err != nil {
		log.Panicln("errore start server: ", err)
	}
}
