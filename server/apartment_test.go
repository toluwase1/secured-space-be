package server

import (
	"encoding/json"
	"fmt"
	"github.com/decadevs/rentals-api/db"
	"github.com/decadevs/rentals-api/models"
	"github.com/decadevs/rentals-api/router"
	"github.com/decadevs/rentals-api/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func Test_CreateApartment(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedDB := db.NewMockDB(ctrl)

	godotenv.Load("../.env")
	s := &Server{
		mockedDB,
		router.NewRouter(),
	}
	route := s.setupRouter()
	user := &models.User{Email: "boboti@gmail.com"}
	apartment := &models.Apartment{
		Models:          models.Models{},
		UserID:          "esdfghjkfhgjh",
		Title:           "Fantastic family flat",
		CategoryID:      "qesrdfhgjhkj",
		Description:     "3 bedroom with awesome neighbors",
		Price:           40000,
		NoOfRooms:       3,
		Furnished:       false,
		Location:        "Lagos",
		ApartmentStatus: false,
		Images:          nil,
		Interiors:       nil,
		Exteriors:       nil,
	}
	marshalledApart, _ := json.Marshal(apartment)

	secret := os.Getenv("JWT_SECRET")
	accessClaims, refreshClaims := services.GenerateClaims(user.Email)
	accToken, _ := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
	services.GenerateToken(&jwt.SigningMethodHMAC{}, refreshClaims, &secret)

	mockedDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false)
	mockedDB.EXPECT().FindUserByEmail(user.Email).Return(user, nil)

	mockedDB.EXPECT().CreateApartment(apartment).Return(apartment, nil)
	t.Run("Testing_For_Apartment_Successfully_Added", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/user/apartments", strings.NewReader(string(marshalledApart)))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		route.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), "Apartment Successfully Added")
	})

	mockedDB.EXPECT().CreateApartment(apartment).Return(nil, errors.New("error creating apartment"))
	t.Run("Testing_For_Error_in_Creating_Apartment", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/user/apartments", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		route.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		assert.Contains(t, rw.Body.String(), "Bad Request")
	})
}
