package server

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/decadevs/rentals-api/db"
	"github.com/decadevs/rentals-api/models"
	"github.com/decadevs/rentals-api/router"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSignupWithCorrectDetailsAgent(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := db.NewMockDB(ctrl)

	s := &Server{
		DB:     m,
		Router: router.NewRouter(),
	}
	r := s.setupRouter()

	user := models.User{
		FirstName: "Spankie",
		LastName:  "Dee",
		Password:  "password",
		Address:   "1, boli drive",
		Email:     "spankie_signup@gmail.com",
		Phone1:    "08909876787",
	}


	m.EXPECT().FindUserByEmail(user.Email).Return(&user, nil)
	jsonuser, err := json.Marshal(user)
	if err != nil {
		t.Fail()
		return
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/signup_agent", strings.NewReader(string(jsonuser)))
	r.ServeHTTP(w, req)


	m.EXPECT().FindUserByEmail(user.Email).Return(&user, nil)
	t.Run("check if user_email exists in the database", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/signup_agent", strings.NewReader(string(jsonuser)))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "User email already exists")
	})

	m.EXPECT().FindUserByEmail(user.Email).Return(&user, errors.New("no record found in database"))
	m.EXPECT().CreateUser(gomock.Any()).Return(nil, nil)
	t.Run("If email does not exist in the database", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/signup_agent", strings.NewReader(string(jsonuser)))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "signup successful")
	})
}