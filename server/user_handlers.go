package server

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/decadevs/rentals-api/models"
	"github.com/decadevs/rentals-api/server/response"
	"github.com/decadevs/rentals-api/servererrors"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// handleShowProfile returns user's details
func (s *Server) handleShowProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		if userI, exists := c.Get("user"); exists {
			if user, ok := userI.(*models.User); ok {
				response.JSON(c, "user details retrieved correctly", http.StatusOK, gin.H{
					"email":      user.Email,
					"phone":      user.Phone,
					"first_name": user.FirstName,
					"last_name":  user.LastName,
					"image":      user.Image,
					"username":   user.Username,
				}, nil)
				return
			}
		}
		log.Printf("can't get user from context\n")
		response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
	}
}

func (s *Server) handleUpdateUserDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		if userI, exists := c.Get("user"); exists {
			if user, ok := userI.(*models.User); ok {

				username, email := user.Username, user.Email
				if errs := s.decode(c, user); errs != nil {
					response.JSON(c, "", http.StatusBadRequest, nil, errs)
					return
				}

				//TODO try to eliminate this
				user.Username, user.Email = username, email
				user.UpdatedAt = time.Now()
				if err := s.DB.UpdateUser(user); err != nil {
					log.Printf("update user error : %v\n", err)
					response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
					return
				}
				response.JSON(c, "user updated successfuly", http.StatusOK, nil, nil)
				return
			}
		}
		log.Printf("can't get user from context\n")
		response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
	}
}

func (s *Server) handleGetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get users should only be access by those have the permission
		if userI, exists := c.Get("user"); exists {
			if user, ok := userI.(*models.User); ok {
				users, err := s.DB.FindAllUsersExcept(user.Email)
				if err != nil {
					log.Printf("find users error : %v\n", err)
					response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
					return
				}
				response.JSON(c, "retrieved users sucessfully", http.StatusOK, gin.H{"users": users}, nil)
				return
			}
		}
		log.Printf("can't get user from context\n")
		response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
		return
	}
}

func (s *Server) handleGetUserByUsername() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := &struct {
			Username string `json:"username" binding:"required"`
		}{}

		if errs := s.decode(c, name); errs != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, errs)
			return
		}

		user, err := s.DB.FindUserByUsername(name.Username)
		if err != nil {
			if inactiveErr, ok := err.(servererrors.InActiveUserError); ok {
				response.JSON(c, "", http.StatusBadRequest, nil, []string{inactiveErr.Error()})
				return
			}
			log.Printf("find user error : %v\n", err)
			response.JSON(c, "user not found", http.StatusNotFound, nil, []string{"user not found"})
			return
		}

		response.JSON(c, "user retrieved successfully", http.StatusOK, gin.H{
			"email":      user.Email,
			"phone":      user.Phone,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"image":      user.Image,
			"username":   user.Username,
		}, nil)
	}
}

// handleUploadProfilePic uploads a user's profile picture
func (s *Server) handleUploadProfilePic() gin.HandlerFunc {
	return func(c *gin.Context) {

		if userI, exists := c.Get("user"); exists {
			if user, ok := userI.(*models.User); ok {

				const maxSize = int64(2048000) // allow only 2MB of file size

				r := c.Request
				err := r.ParseMultipartForm(maxSize)
				if err != nil {
					log.Printf("parse image error: %v\n", err)
					response.JSON(c, "", http.StatusBadRequest, nil, []string{"image too large"})
					return
				}

				file, fileHeader, err := r.FormFile("profile_picture")
				if err != nil {
					log.Println(err)
					response.JSON(c, "", http.StatusBadRequest, nil, []string{"image not supplied"})
					return
				}
				defer file.Close()

				supportedFileTypes := map[string]bool{
					".png":  true,
					".jpeg": true,
					".jpg":  true,
				}
				fileExtension := filepath.Ext(fileHeader.Filename)
				if !supportedFileTypes[fileExtension] {
					log.Println(fileExtension)
					response.JSON(c, "", http.StatusBadRequest, nil, []string{fileExtension + " image file type is not supported"})
					return
				}
				tempFileName := "profile_pics/" + bson.NewObjectId().Hex() + fileExtension

				session, err := session.NewSession(&aws.Config{
					Region: aws.String(os.Getenv("AWS_REGION")),
					Credentials: credentials.NewStaticCredentials(
						os.Getenv("AWS_SECRET_ID"),
						os.Getenv("AWS_SECRET_KEY"),
						os.Getenv("AWS_TOKEN"),
					),
				})
				if err != nil {
					log.Printf("could not upload file: %v\n", err)
				}

				err = uploadFileToS3(session, file, tempFileName, fileHeader.Size)
				if err != nil {
					log.Println(err)
					response.JSON(c, "", http.StatusInternalServerError, nil, []string{"an error occured while uploading the image"})
					return
				}

				user.Image = os.Getenv("S3_BUCKET") + tempFileName
				if err = s.DB.UpdateUser(user); err != nil {
					log.Println(err)
					response.JSON(c, "", http.StatusInternalServerError, nil, []string{"unable to update user's profile pic"})
					return
				}

				response.JSON(c, "successfully created file", http.StatusOK, gin.H{
					"imageurl": user.Image,
				}, nil)
				return
			}
		}
		response.JSON(c, "", http.StatusUnauthorized, nil, []string{"unable to retrieve authenticated user"})
		return
	}
}
