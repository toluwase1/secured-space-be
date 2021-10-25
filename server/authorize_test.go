package server

import (
	"fmt"
	"github.com/decadevs/rentals-api/db"
	"github.com/decadevs/rentals-api/router"
	"github.com/decadevs/rentals-api/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)
func TestUnAuthorize(t *testing.T) {
	if err := godotenv.Load(); err != nil {
		t.Fail()
	}
	ctrl := gomock.NewController(t)
	m := db.NewMockDB(ctrl)

	s := &Server{
		DB:     m,
		Router: router.NewRouter(),
	}
	router := s.setupRouter()

	t.Run("Test_For_Wrong_Secret_Key", func(t *testing.T) {
		accessClaims, _ := services.GenerateClaims("adebayo@gmail.com")

		secret := "jakhjb67273hsbv"
		accToken, err := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
		if err != nil {
			t.Fail()
		}
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer-%s", *accToken))

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
	})

	t.Run("Test_For_Blacklist", func(t *testing.T) {
		accessClaims, _ := services.GenerateClaims("adebayo@gmail.com")

		secret := os.Getenv("JWT_SECRET")
		accToken, err := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
		if err != nil {
			t.Fail()
		}
		m.EXPECT().TokenInBlacklist(accToken).Return(true)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer-%s", *accToken))

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
	})
}

//func TestAuthorize(t *testing.T) {
//	if err := godotenv.Load(); err != nil {
//		t.Fail()
//	}
//	ctrl := gomock.NewController(t)
//	m := db.NewMockDB(ctrl)
//
//	s := &Server{
//		DB:     m,
//		Router: router.NewRouter(),
//	}
//	router := s.setupRouter()
//
//	t.Run("Test_Authorize", func(t *testing.T) {
//		accessClaims, _ := services.GenerateClaims("adebayo@gmail.com")
//
//		secret := os.Getenv("JWT_SECRET")
//		accToken, err := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
//		if err != nil {
//			t.Fail()
//		}
//
//		user := &models.User{Email: "adebayo@gmail.com"}
//		users := []models.User{{}}
//		m.EXPECT().FindUserByEmail(user.Email).Return(user, nil)
//		m.EXPECT().TokenInBlacklist(accToken).Return(false)
//		m.EXPECT().FindAllUsersExcept(user.Email).Return(users, nil)
//
//		w := httptest.NewRecorder()
//		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
//		req.Header.Set("Authorization", fmt.Sprintf("Bearer-%s", *accToken))
//
//		router.ServeHTTP(w, req)
//
//		assert.Equal(t, http.StatusOK, w.Code)
//		assert.Contains(t, w.Body.String(), "retrieved users successfully")
//	})
//
//	t.Run("Test_FIndUserByEmail_Error", func(t *testing.T) {
//		accessClaims, _ := services.GenerateClaims("adebayo@gmail.com")
//
//		secret := os.Getenv("JWT_SECRET")
//		accToken, err := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
//		if err != nil {
//			t.Fail()
//		}
//
//		user := &models.User{Email: "adebayo@gmail.com"}
//		m.EXPECT().FindUserByEmail(user.Email).Return(user, errors.New("an error occurred"))
//		m.EXPECT().TokenInBlacklist(accToken).Return(false)
//
//		w := httptest.NewRecorder()
//		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
//		req.Header.Set("Authorization", fmt.Sprintf("Bearer-%s", *accToken))
//
//		router.ServeHTTP(w, req)
//		assert.Equal(t, http.StatusNotFound, w.Code)
//		assert.Contains(t, w.Body.String(), "user not found")
//	})
//}
