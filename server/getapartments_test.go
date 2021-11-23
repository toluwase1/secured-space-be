package server

import (
	"fmt"
	"github.com/decadevs/rentals-api/db"
	"github.com/decadevs/rentals-api/models"
	"github.com/decadevs/rentals-api/router"
	"github.com/decadevs/rentals-api/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func AuthorizeTestRoute(mDB *db.MockDB, t *testing.T, email string) (*models.User, *string) {
	accessClaims, refreshClaims := services.GenerateClaims(email)

	secret := os.Getenv("JWT_SECRET")
	accToken, err := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
	if err != nil {
		t.Fail()
	}
	services.GenerateToken(jwt.SigningMethodHS256, refreshClaims, &secret)

	user := &models.User{Email: email}
	user.ID = "123456rtgfdvdsawer"
	mDB.EXPECT().FindUserByEmail(user.Email).Return(user, nil)
	mDB.EXPECT().TokenInBlacklist(accToken).Return(false)
	return user, accToken
}

func TestApplication_GetUserApartments(t *testing.T) {
	ctrl := gomock.NewController(t)
	mDB := db.NewMockDB(ctrl)

	s := &Server{
		DB:     mDB,
		Router: router.NewRouter(),
	}
	router := s.setupRouter()

	t.Run("Testing_Error_Getting_All_Apartments", func(t *testing.T) {
		user, accToken := AuthorizeTestRoute(mDB, t, "tchisom@gmail.com")
		mDB.EXPECT().GetUsersApartments(user.ID).Return(nil, errors.New("an error occurred"))
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/user-apartment", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		router.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		assert.Contains(t, rw.Body.String(), "internal server error")
	})

	t.Run("Test_For_Successful_Response", func(t *testing.T) {
		user, accToken := AuthorizeTestRoute(mDB, t, "tchisom@gmail.com")
		mDB.EXPECT().GetUsersApartments(user.ID).Return(nil, nil)
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/user-apartment", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		router.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), "retrieved apartments successfully")
	})
}
