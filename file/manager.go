package file

import (
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/h2non/filetype"
)

// ListFile ...
func ListFile() []os.FileInfo {
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

// SaveFile ...
func SaveFile(part *multipart.Part) error {
	// prepare the dst
	dst, err := os.Create("./data/" + part.FileName())
	defer dst.Close()
	if err != nil {
		return err
	}
	// copy the part to dst
	if _, err := io.Copy(dst, part); err != nil {
		return err
	}
	return nil
}

// GetFile ...
func GetFile(name string) (*os.File, error) {
	file, err := os.Open("./data/" + name)
	return file, err
}

// MimetypeFile ...
func MimetypeFile(file *os.File) (string, error) {
	// To calculate the mimetype we only have to pass the file header = first 261 bytes
	head := make([]byte, 261)
	file.Read(head)
	mimetype, err := filetype.Get(head)
	if err != nil {
		log.Printf("Error detect mimetype: %v", err)
		return "", err
	}
	// prepare response, set the offset to zero for the next Read
	file.Seek(0, 0)
	return mimetype.MIME.Value, nil
}
