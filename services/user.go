package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/globalsign/mgo/bson"
	"os"
	"path/filepath"
)

func CheckSupportedFile(filename string) (string, bool) {
	supportedFileTypes := map[string]bool{
		".png":  true,
		".jpeg": true,
		".jpg":  true,
	}
	fileExtension := filepath.Ext(filename)

	return fileExtension, !supportedFileTypes[fileExtension]
}

func PreAWS(fileExtension, folder string) (*session.Session, string, error) {
	tempFileName := folder + "/" + bson.NewObjectId().Hex() + fileExtension

	session, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_SECRET_ID"),
			os.Getenv("AWS_SECRET_KEY"),
			os.Getenv("AWS_TOKEN"),
		),
	})

	return session, tempFileName, err
}
