package server

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/decadevs/rentals-api/models"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"mime/multipart"
	"net/http"
	"os"
)

func uploadFileToS3(s *session.Session, file multipart.File, fileName string, size int64) error {
	// get the file size and read
	// the file content into a buffer
	buffer := make([]byte, size)
	file.Read(buffer)
	// config settings: this is where you choose the bucket,
	// filename, content-type and storage class of the file
	// you're uploading
	_, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:                  aws.String(fileName),
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(int64(size)),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("INTELLIGENT_TIERING"),
	})
	return err
}

func GetUserFromContext(c *gin.Context) (*models.User, error){
	if userI, exists := c.Get("user"); exists {
		if user, ok := userI.(*models.User); ok {
			return user, nil
		}
		return nil, errors.New("User is not logged in")
	}
	return nil, errors.New("user is not logged in")

}

