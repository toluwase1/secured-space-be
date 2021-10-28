package server

import (
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
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestServer_DeleteApartment(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println(err.Error())
	}
	ctrl := gomock.NewController(t)
	mockDB := db.NewMockDB(ctrl)
	s := &Server{
		DB: mockDB,
		Router: router.NewRouter(),
	}
	r := s.setupRouter()

	accessClaims, _ := services.GenerateClaims("adebayo@gmail.com")

	secret := os.Getenv("JWT_SECRET")
	accToken, err := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
	if err != nil {
		t.Fail()
	}
	user := &models.User{Email: "adebayo@gmail.com"}
	user.ID = "aefrfh123435waes"
	mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false).Times(3)
	mockDB.EXPECT().FindUserByEmail(user.Email).Return(user, nil).Times(3)
	t.Run("Test_Empty_ApartmentID", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/api/v1/user/apartment//", nil)
		if err != nil {
			t.Fail()
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		assert.Contains(t, rw.Body.String(),"apartment id cannot be empty")
	})
	mockDB.EXPECT().DeleteApartment("12323shjbvbhj1t", user.ID).Return(errors.New("an error occurred"))
	t.Run("Test_Error-Deleting_Apartment", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/api/v1/user/apartment/12323shjbvbhj1t/", nil)
		if err != nil {
			t.Fail()
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		assert.Contains(t, rw.Body.String(), "internal server error")
	})

	mockDB.EXPECT().DeleteApartment("12323shjbvbhj1t", user.ID).Return(nil)
	t.Run("Test_For_Successful_Delete", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/api/v1/user/apartment/12323shjbvbhj1t/", nil)
		if err != nil {
			t.Fail()
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), "Deleted Successfully")
	})
}
