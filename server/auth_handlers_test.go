package server

import (
	"encoding/json"
	"fmt"
	"github.com/decadevs/rentals-api/db"
	"github.com/decadevs/rentals-api/models"
	router2 "github.com/decadevs/rentals-api/router"
	"github.com/decadevs/rentals-api/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestApplication_HandleLogout(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := db.NewMockDB(ctrl)

	s := &Server{
		DB:     m,
		Router: router2.NewRouter(),
	}
	router := s.setupRouter()
	user := &models.User{
		Email: "yemmy@gmail.com",
	}
	user.ID = "123456ygszxyut54"
	accessClaims, refreshClaims := services.GenerateClaims("yemmy@gmail.com")

	secret := os.Getenv("JWT_SECRET")
	accToken, err := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
	refToken, err := services.GenerateToken(jwt.SigningMethodHS256, refreshClaims, &secret)
	if err != nil {
		t.Fail()
	}
	m.EXPECT().FindUserByEmail(user.Email).Return(user, nil)
	m.EXPECT().TokenInBlacklist(accToken).Return(false)

	rt := &struct {
		RefreshToken string `json:"refresh_token,omitempty" binding:"required"`
	}{
		RefreshToken: *refToken,
	}
	marshalledRefreshToken, _ := json.Marshal(rt)

	m.EXPECT().AddToBlackList(gomock.Any()).Return(nil)

	m.EXPECT().AddToBlackList(gomock.Any()).Return(nil)
	t.Run("Test_Logout", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/logout", strings.NewReader(string(marshalledRefreshToken)))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))

		router.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), "logout successful")
	})
}

func Test_handleLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedDB := db.NewMockDB(ctrl)

	s := &Server{
		DB: mockedDB,
	}

	router := s.setupRouter()

	t.Run("Test_For_Login_Request", func(t *testing.T) {
		loginRequest := &struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}{
			Email:    "jdoe@gmail.com",
			Password: "",
		}
		jsonFile, err := json.Marshal(loginRequest)
		if err != nil {
			t.Error("Failed to marshal file")
		}
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(string(jsonFile)))

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "validation failed on field 'Password', condition: required")
	})

	t.Run("Test_FindUserByEmail", func(t *testing.T) {
		hashedP, _ := services.GenerateHashPassword("password")
		user := &models.User{Email: "jdoe@gmail.com", HashedPassword: string(hashedP)}
		mockedDB.EXPECT().FindUserByEmail(user.Email).Return(user, nil)
		loginRequest := &struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}{
			Email:    "jdoe@gmail.com",
			Password: "password",
		}
		jsonFile, err := json.Marshal(loginRequest)
		if err != nil {
			t.Error("Failed to marshal file")
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(string(jsonFile)))

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "")
	})
}
