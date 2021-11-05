package server

import (
	"encoding/json"
	"github.com/decadevs/rentals-api/db"
	"github.com/decadevs/rentals-api/models"
	"github.com/decadevs/rentals-api/services"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
