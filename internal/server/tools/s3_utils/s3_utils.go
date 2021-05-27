package s3_utils

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/configs"
	_ "github.com/go-park-mail-ru/2021_1_DuckLuck/configs"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
)

var (
	sess     *session.Session
	bucket   string
	acl      string
	endpoint string
	svc      *s3.S3
)

func InitS3() {
	// Load s3 environment
	err := godotenv.Load(configs.ApiServerS3Env)
	if err != nil {
		log.Fatal(err)
	}

	accessKeyID := os.Getenv("S3_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("S3_SECRET_ACCESS_KEY")
	myRegion := os.Getenv("S3_REGION")
	bucket = os.Getenv("S3_BUCKET")
	acl = os.Getenv("S3_ACL")
	endpoint = os.Getenv("S3_ENDPOINT")
	//var err error
	sess, err = session.NewSession(
		&aws.Config{
			Region:   aws.String(myRegion),
			Endpoint: aws.String(endpoint),
			Credentials: credentials.NewStaticCredentials(
				accessKeyID,
				secretAccessKey,
				"",
			),
		})
	svc = s3.New(sess)
	if err != nil {
		log.Fatal(err)
	}
}

func UploadMultipartFile(path string, file *multipart.File, header *multipart.FileHeader) (string, error) {
	var fileType string
	switch header.Header.Get("Content-Type") {
	case "image/png":
		fileType = ".png"
	case "image/jpg":
		fileType = ".jpg"
	case "image/jpeg":
		fileType = ".jpeg"
	default:
		return "", errors.ErrIncorrectFileType
	}

	newName := path + "/" + uuid.NewV4().String() + fileType
	uploader := s3manager.NewUploader(sess)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		ACL:    aws.String(acl),
		Key:    aws.String(newName),
		Body:   *file,
	})
	if err != nil {
		return "", errors.ErrS3InternalError
	}

	return newName, nil
}

func DeleteFile(fileName string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	}

	_, err := svc.DeleteObject(input)
	if err != nil {
		return errors.ErrS3InternalError
	}

	return nil
}

func PathToFile(fileName string) string {
	return fmt.Sprintf("https://%s.%s/%s", bucket, endpoint, fileName)
}
