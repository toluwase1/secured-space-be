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

//func HandleFileUpload(file *multipart.File, header multipart.FileHeader) (string,error){
//	fileExtension, ok := CheckSupportedFile(strings.ToLower(header.Filename))
//	log.Printf(filepath.Ext(strings.ToLower(header.Filename)))
//	fmt.Println(fileExtension)
//	if ok {
//		log.Println(fileExtension)
//		response.JSON(c, "", http.StatusBadRequest, nil, []string{fileExtension + " image file type is not supported"})
//		return
//	}
//
//	session, tempFileName, err := PreAWS(fileExtension)
//
//	if err != nil {
//		log.Printf("could not upload file: %v\n", err)
//	}
//
//	err = DB.UploadFileToS3(session, file, tempFileName, fileHeader.Size)
//	if err != nil {
//		log.Println(err)
//		response.JSON(c, "", http.StatusInternalServerError, nil, []string{"an error occured while uploading the image"})
//		return
//	}
//	return "", nil
//}
