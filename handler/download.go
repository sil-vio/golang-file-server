package handler

import (
	"fmt"
	fileManager "github.com/sil-vio/golang-file-server/file"
	"io"
	"log"
	"net/http"
)

// DownloadHandler ..
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// extract from url the query param filename
		queryParam := r.URL.Query()
		filename := queryParam.Get("filename")
		if filename == "" {
			log.Printf("No filename specified")
			http.Error(w, "No filename specified", http.StatusBadRequest)
			return
		}
		// get the requested file
		file, err := fileManager.GetFile(filename)
		if err != nil {
			log.Printf("can't find requested file: %v", err)
			http.Error(w, "can't find requested file", http.StatusBadRequest)
			return
		}
		mimetype, err := fileManager.MimetypeFile(file)
		if err != nil {
			http.Error(w, "Error detect mimetype", http.StatusBadRequest)
		}
		contentDisposition := "attachment; filename=" + filename
		w.Header().Set("Content-Disposition", contentDisposition)
		w.Header().Set("Content-Type", mimetype)
		info, _ := file.Stat()
		w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))

		//stream the body to the client without fully loading it into memory
		io.Copy(w, file)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
