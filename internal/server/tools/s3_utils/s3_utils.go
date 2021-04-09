package file_utils

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/configs"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	uuid "github.com/satori/go.uuid"
)

var (
	sessionS3 *session.Session
	bucketS3 string
	aclS3 string
)

func init() {
	// Load environment
	err := godotenv.Load(configs.PathToConfigEnv)
	if err != nil {
		log.Fatal(err)
	}

	accessKeyID := os.Getenv("S3_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("S3_SECRET_ACCESS_KEY")
	myRegion := os.Getenv("S3_REGION")
	bucketS3 = os.Getenv("S3_BUCKET")
	aclS3 = os.Getenv("S3_ACL")
	//var err error
	sessionS3, err = session.NewSession(
		&aws.Config{
			Region:   aws.String(myRegion),
			Endpoint: aws.String(os.Getenv("S3_ENDPOINT")),
			Credentials: credentials.NewStaticCredentials(
				accessKeyID,
				secretAccessKey,
				"",
			),
		})
	if err != nil {
		log.Fatal(err)
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
	res, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketS3),
		ACL:    aws.String(aclS3),
		Key:    aws.String(newName),
		Body:   file,
	})
	if err != nil {
		return "", errors.ErrUploadToS3
	}

	return res.Location, nil
}
