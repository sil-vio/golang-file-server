package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/h2non/filetype"
)

// compiling/caching the template

var templates = template.Must(template.New("index.html").ParseFiles("./templates/index.html"))

func listFile() []os.FileInfo {
	f, err := os.Open("./data/")
	if err != nil {
		log.Print(err)
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Print(err)
	}
	onlyfiles := make([]os.FileInfo, 0)

	for _, file := range files {
		if !file.IsDir() {
			onlyfiles = append(onlyfiles, file)
		}
	}
	return onlyfiles
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	// GET to display the upload form.
	case "GET":

		err := templates.Execute(w, map[string]interface{}{
			"Files": listFile()})
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

			// prepare the dst
			dst, err := os.Create("./data/" + part.FileName())
			defer dst.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// copy the part to dst
			if _, err := io.Copy(dst, part); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		// displaying a success message.
		err = templates.Execute(w, map[string]interface{}{
			"Files": listFile(), "Message": "Upload successful."})
		if err != nil {
			log.Print(err)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
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
		// open the requested file
		log.Println("File requested: ", filename)
		file, err := os.Open("./data/" + filename)
		if err != nil {
			log.Printf("can't find requested file: %v", err)
			http.Error(w, "can't find requested file", http.StatusBadRequest)
			return
		}
		// To calculate the mimetype we only have to pass the file header = first 261 bytes
		head := make([]byte, 261)
		file.Read(head)
		mimetype, err := filetype.Get(head)
		if err != nil {
			log.Printf("Error detect mimetype: %v", err)
			http.Error(w, "Error detect mimetype", http.StatusBadRequest)
			return
		}
		// prepare response, set the offset to zero for the next Read
		file.Seek(0, 0)
		contentDisposition := "attachment; filename=" + filename
		w.Header().Set("Content-Disposition", contentDisposition)
		w.Header().Set("Content-Type", mimetype.MIME.Value)
		info, _ := file.Stat()
		w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))

		//stream the body to the client without fully loading it into memory
		io.Copy(w, file)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/site", uploadHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)
	log.Print("Listening on port:4200...")
	// Listen on port 8080
	err := http.ListenAndServe(":4200", nil)
	if err != nil {
		log.Panicln("errore start server: ", err)
	}
}
