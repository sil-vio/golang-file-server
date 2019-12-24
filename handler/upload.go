package handler

import (
	fileManager "github.com/sil-vio/golang-file-server/file"
	"html/template"
	"io"
	"log"
	"net/http"
)

// compiling/caching the template
var templates = template.Must(template.New("index.html").ParseFiles("./templates/index.html"))

// UploadHandler ..
func UploadHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	// GET to display the upload form.
	case "GET":

		err := templates.Execute(w, map[string]interface{}{
			"Files": fileManager.ListFile()})
		if err != nil {
			log.Print(err)
		}
		// POST analyzes each part of the MultiPartReader (ie the uploaded file(s))
		// and saves them to disk.
	case "POST":
		// grab the request.MultipartReader
		reader, err := r.MultipartReader()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// copy each part to destination.
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			// if part.FileName() is empty, skip this iteration.
			if part.FileName() == "" {
				continue
			}
			fileManager.SaveFile(part)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		// displaying a success message.
		err = templates.Execute(w, map[string]interface{}{
			"Files": fileManager.ListFile(), "Message": "Upload successful."})
		if err != nil {
			log.Print(err)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
