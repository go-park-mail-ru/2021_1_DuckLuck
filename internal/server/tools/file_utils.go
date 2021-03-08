package tools

import (
	"io"
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	uuid "github.com/satori/go.uuid"
)

func UploadFile(r *http.Request, fileKey string, pathToUpload string) (string, error) {
	// Max size - 10 Mb
	r.ParseMultipartForm(10 * 1024 * 1024)
	file, handler, err := r.FormFile(fileKey)
	if err != nil {
		return "", errors.ErrFileNotRead
	}
	defer file.Close()

	var fileType string
	switch handler.Header.Get("Content-Type") {
	case "image/png":
		fileType = ".png"
	case "image/jpg":
		fileType = ".jpg"
	case "image/jpeg":
		fileType = ".jpeg"
	default:
		return "", errors.ErrIncorrectFileType
	}

	newName := uuid.NewV4().String() + fileType
	newFile, err := os.Create(pathToUpload + newName)
	if err != nil {
		return "", errors.ErrServerSystem
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, file)
	if err != nil {
		return "", errors.ErrServerSystem
	}

	return newName, nil
}
