package server

import (
	"encoding/json"
	"fmt"
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

func TestSignupWithInCorrectDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := db.NewMockDB(ctrl)

	s := &Server{
		DB:     m,
		Router: router.NewRouter(),
	}
	router := s.setupRouter()

	user := models.User{
		FirstName: "Spankie",
		LastName:  "Dee",
		Password:  "password",
		Email:     "spankie_signup",
		Phone1:    "08909876787",
	}

	jsonuser, err := json.Marshal(user)
	if err != nil {
		t.Fail()
		return
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/signup_tenant", strings.NewReader(string(jsonuser)))
	router.ServeHTTP(w, req)

	bodyString := w.Body.String()
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, bodyString, fmt.Sprintf("validation failed on field 'Email', condition: email, actual: %s", user.Email))
}

func TestSignupIfEmailExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := db.NewMockDB(ctrl)

	s := &Server{
		DB:     m,
		Router: router.NewRouter(),
	}
	router := s.setupRouter()

	user := models.User{
		FirstName: "Spankie",
		LastName:  "Dee",
		Password:  "password",
		Address:   "1, boli drive",
		Email:     "spankie_signup@gmail.com",
		Phone1:    "08909876787",
	}

	m.EXPECT().FindUserByEmail(gomock.Any()).Return(&user, nil)

	jsonuser, err := json.Marshal(user)
	if err != nil {
		t.Fail()
		return
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/signup_tenant", strings.NewReader(string(jsonuser)))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "User email already exists")
}

func TestSignupWithCorrectDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := db.NewMockDB(ctrl)

	s := &Server{
		DB:     m,
		Router: router.NewRouter(),
	}
	router := s.setupRouter()

	user := models.User{
		FirstName: "Spankie",
		LastName:  "Dee",
		Password:  "password",
		Address:   "1, boli drive",
		Email:     "spankie_signup@gmail.com",
		Phone1:    "08909876787",
	}

	m.EXPECT().CreateUser(user.Email).Return(&user, nil)

	jsonuser, err := json.Marshal(user)
	if err != nil {
		t.Fail()
		return
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/signup_tenant", strings.NewReader(string(jsonuser)))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "signup successful")
}
