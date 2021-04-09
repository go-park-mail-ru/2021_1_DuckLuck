package file_utils

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	uuid "github.com/satori/go.uuid"
)

var sessionS3 *session.Session

func init() {
	accessKeyID := "3rnonniNuYMrbnrr519hes"
	secretAccessKey := "an9weAagS9F3N8L2J2zned3wEgec5v438sMdeqy13KMz"
	myRegion := "ap-northeast-1"
	var err error
	sessionS3, err = session.NewSession(
		&aws.Config{
			Region:   aws.String(myRegion),
			Endpoint: aws.String("hb.bizmrg.com"),
			Credentials: credentials.NewStaticCredentials(
				accessKeyID,
				secretAccessKey,
				"",
			),
		})
	if err != nil {
		panic(err)
	}
}

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
	uploader := s3manager.NewUploader(sessionS3)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("duckluckmarket"),
		ACL:    aws.String("public-read"),
		Key:    aws.String("test_file_name"),
		Body:   file,
	})
	fmt.Println(err)

	return newName, nil
}
