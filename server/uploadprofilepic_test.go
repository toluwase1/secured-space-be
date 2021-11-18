//package server
//
//import (
//	"bytes"
//	"fmt"
//	"github.com/decadevs/rentals-api/db"
//	"github.com/decadevs/rentals-api/models"
//	router2 "github.com/decadevs/rentals-api/router"
//	"github.com/decadevs/rentals-api/services"
//	"github.com/dgrijalva/jwt-go"
//	"github.com/golang/mock/gomock"
//	"github.com/stretchr/testify/assert"
//	"image"
//	"image/color"
//	"image/png"
//	"mime/multipart"
//	"net/http"
//	"net/http/httptest"
//	"os"
//	"testing"
//)
//
//func createImg() (*os.File, error) {
//	width := 200
//	height := 100
//
//	upLeft := image.Point{0, 0}
//	lowRight := image.Point{width, height}
//
//	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
//
//	// Colors are defined by Red, Green, Blue, Alpha uint8 values.
//	cyan := color.RGBA{100, 200, 200, 0xff}
//
//	// Set color for each pixel.
//	for x := 0; x < width; x++ {
//		for y := 0; y < height; y++ {
//			switch {
//			case x < width/2 && y < height/2: // upper left quadrant
//				img.Set(x, y, cyan)
//			case x >= width/2 && y >= height/2: // lower right quadrant
//				img.Set(x, y, color.White)
//			default:
//				// Use zero value.
//			}
//		}
//	}
//
//	// Encode as PNG.
//	f, err := os.Create("testimg.png")
//	if err != nil {
//		return nil, err
//	}
//	png.Encode(f, img)
//	return f, nil
//}
//func prepfile(file *os.File) (*bytes.Buffer, string, error) {
//	var b bytes.Buffer
//	w := multipart.NewWriter(&b)
//	defer w.Close()
//	fmt.Println(file.Name())
//	if _, err := w.CreateFormFile("profile_picture", file.Name()); err != nil {
//		return nil, "", fmt.Errorf("%v", err)
//	}
//	fmt.Println(w.FormDataContentType())
//	return &b, w.FormDataContentType(), nil
//}
//func TestUploadprofielpic(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	m := db.NewMockDB(ctrl)
//
//	s := &Server{
//		DB:     m,
//		Router: router2.NewRouter(),
//	}
//	accessClaims, _ := services.GenerateClaims("chinonso@gmail.com")
//	secret := os.Getenv("JWT_SECRET")
//	accToken, err := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
//	if err != nil {
//		t.Fail()
//	}
//	file, err := createImg()
//	if err != nil {
//		panic(err)
//	}
//
//	user := &models.User{Email: "chinonso@gmail.com"}
//	m.EXPECT().FindUserByEmail(user.Email).Return(user, nil)
//	m.EXPECT().TokenInBlacklist(accToken).Return(false)
//	m.EXPECT().UploadFileToS3(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
//
//	router := s.setupRouter()
//
//	t.Run("TestUploadprofielpic", func(t *testing.T) {
//
//		b, content_type, err := prepfile(file)
//		if err != nil {
//			fmt.Errorf("%v", err)
//		}
//
//		resp := httptest.NewRecorder()
//
//		req, _ := http.NewRequest(http.MethodPost, "/api/v1/me/uploadpic", b)
//		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
//		req.Header.Set("Content-Type", content_type)
//
//		router.ServeHTTP(resp, req)
//		assert.Equal(t, http.StatusOK, 200)
//
//		//This how far i can go and i don't know how to set the key for file to so that it can be read from the body
//	})
//}
