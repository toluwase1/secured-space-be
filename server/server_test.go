package server

import (
	"encoding/json"
	"fmt"
	"github.com/decadevs/rentals-api/db"
	"github.com/decadevs/rentals-api/models"
	"github.com/decadevs/rentals-api/router"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSignupWithInCorrectDetailsTenant(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := db.NewMockDB(ctrl)

	s := &Server{
		DB:     m,
		Router: router.NewRouter(),
	}
	r := s.setupRouter()



	role := models.Role{
		Models: models.Models{},
		Title:  "tenant",
	}

	user := models.User{
		FirstName: "Spankie",
		LastName:  "Dee",
		Password:  "password",
		Email:     "spankie_signup",
		Phone1:    "08909876787",
		RoleID: role.ID,
		Role: role,
	}
	m.EXPECT().GetRoleByName("tenant").Return(role, nil)

	jsonuser, err := json.Marshal(user)
	if err != nil {
		t.Fail()
		return
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/signup_tenant", strings.NewReader(string(jsonuser)))
	r.ServeHTTP(w, req)

	bodyString := w.Body.String()
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, bodyString, fmt.Sprintf("validation failed on field 'Email', condition: email, actual: %s", user.Email))
}

func TestSignupIfEmailExistsTenant(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := db.NewMockDB(ctrl)

	s := &Server{
		DB:     m,
		Router: router.NewRouter(),
	}
	r := s.setupRouter()

	role := models.Role{
		Models: models.Models{},
		Title:  "tenant",
	}

	user := models.User{
		FirstName: "Spankie",
		LastName:  "Dee",
		Password:  "password",
		Address:   "1, boli drive",
		Email:     "spankie_signup@gmail.com",
		Phone1:    "08909876787",
		RoleID: role.ID,
		Role: role,
	}

	m.EXPECT().GetRoleByName("tenant").Return(role, nil)
	m.EXPECT().FindUserByEmail(gomock.Any()).Return(&user, nil)

	jsonuser, err := json.Marshal(user)
	if err != nil {
		t.Fail()
		return
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/signup_tenant", strings.NewReader(string(jsonuser)))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "User email already exists")
}

func TestSignupWithCorrectDetailsTenant(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := db.NewMockDB(ctrl)

	s := &Server{
		DB:     m,
		Router: router.NewRouter(),
	}
	r := s.setupRouter()


	role := models.Role{
		Models: models.Models{},
		Title:  "tenant",
	}

	user := models.User{
		FirstName: "Spankie",
		LastName:  "Dee",
		Password:  "password",
		Address:   "1, boli drive",
		Email:     "spankie_signup@gmail.com",
		Phone1:    "08909876787",
		RoleID: role.ID,
		Role: role,
	}

	m.EXPECT().GetRoleByName("tenant").Return(role, nil)
	m.EXPECT().FindUserByEmail(user.Email).Return(&user, nil)
	jsonuser, err := json.Marshal(user)
	if err != nil {
		t.Fail()
		return
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/signup_tenant", strings.NewReader(string(jsonuser)))
	r.ServeHTTP(w, req)

	m.EXPECT().GetRoleByName("tenant").Return(role, nil)
	m.EXPECT().FindUserByEmail(user.Email).Return(&user, nil)
	t.Run("check if tenant_email exists in the database", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/signup_tenant", strings.NewReader(string(jsonuser)))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "User email already exists")
	})

}

func TestSignupWithInCorrectDetailsAgent(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := db.NewMockDB(ctrl)

	s := &Server{
		DB:     m,
		Router: router.NewRouter(),
	}
	r := s.setupRouter()

	role := models.Role{
		Models: models.Models{},
		Title:  "agent",
	}

	user := models.User{
		FirstName: "Spankie",
		LastName:  "Dee",
		Password:  "password",
		Email:     "spankie_signup",
		Phone1:    "08909876787",
		RoleID: role.ID,
		Role: role,
	}

	m.EXPECT().GetRoleByName("agent").Return(role, nil)


	jsonuser, err := json.Marshal(user)
	if err != nil {
		t.Fail()
		return
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/signup_agent", strings.NewReader(string(jsonuser)))
	r.ServeHTTP(w, req)

	bodyString := w.Body.String()
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, bodyString, fmt.Sprintf("validation failed on field 'Email', condition: email, actual: %s", user.Email))
}

func TestSignupIfEmailExistsAgent(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := db.NewMockDB(ctrl)

	s := &Server{
		DB:     m,
		Router: router.NewRouter(),
	}
	r := s.setupRouter()

	role := models.Role{
		Models: models.Models{},
		Title:  "agent",
	}

	user := models.User{
		FirstName: "Spankie",
		LastName:  "Dee",
		Password:  "password",
		Address:   "1, boli drive",
		Email:     "spankie_signup@gmail.com",
		Phone1:    "08909876787",
		RoleID: role.ID,
		Role: role,
	}

	m.EXPECT().GetRoleByName("agent").Return(role, nil)
	m.EXPECT().FindUserByEmail(gomock.Any()).Return(&user, nil)

	jsonuser, err := json.Marshal(user)
	if err != nil {
		t.Fail()
		return
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/signup_agent", strings.NewReader(string(jsonuser)))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "User email already exists")
}

func TestSignupWithCorrectDetailsAgent(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := db.NewMockDB(ctrl)

	s := &Server{
		DB:     m,
		Router: router.NewRouter(),
	}
	r := s.setupRouter()

	role := models.Role{
		Models: models.Models{},
		Title:  "agent",
	}

	user := models.User{
		FirstName: "Spankie",
		LastName:  "Dee",
		Password:  "password",
		Address:   "1, boli drive",
		Email:     "spankie_signup@gmail.com",
		Phone1:    "08909876787",
		RoleID: role.ID,
		Role: role,
	}

	m.EXPECT().GetRoleByName("agent").Return(role, nil)
	m.EXPECT().FindUserByEmail(user.Email).Return(&user, nil)
	jsonuser, err := json.Marshal(user)
	if err != nil {
		t.Fail()
		return
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/signup_agent", strings.NewReader(string(jsonuser)))
	r.ServeHTTP(w, req)

	m.EXPECT().GetRoleByName("agent").Return(role, nil)
	m.EXPECT().FindUserByEmail(user.Email).Return(&user, nil)
	t.Run("check if user_email exists in the database", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/signup_agent", strings.NewReader(string(jsonuser)))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "User email already exists")
	})


}
