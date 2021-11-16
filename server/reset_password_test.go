package server

import (
	"encoding/json"
	"github.com/decadevs/rentals-api/db"
	"github.com/decadevs/rentals-api/router"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestApplication_ResetPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedDb := db.NewMockDB(ctrl)

	s := &Server{
		DB:     mockedDb,
		Router: router.NewRouter(),
	}
	r := s.setupRouter()

	t.Run("Test_For_Missing_Password", func(t *testing.T) {
		requestPassword := struct {
			Password string
		}{
			Password: "",
		}
		body, _ := json.Marshal(requestPassword)
		req, err := http.NewRequest(http.MethodPost, "/api/v1/reset-password/12340987wertt", strings.NewReader(string(body)))
		if err != nil {
			t.Errorf("Error Creating Request: %v", err.Error())
		}
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)
		assert.Contains(t, res.Body.String(), "validation failed on field 'Password'")
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})
	requestPassword := struct {
		Password string
	}{
		Password: "Password",
	}
	mockedDb.EXPECT().ResetPassword("12340987wertt", gomock.Any()).Return(errors.New("an error occurred"))
	t.Run("Test_For_Reset_Password_Error", func(t *testing.T) {
		body, _ := json.Marshal(requestPassword)
		req, err := http.NewRequest(http.MethodPost, "/api/v1/reset-password/12340987wertt", strings.NewReader(string(body)))
		if err != nil {
			t.Errorf("Error Creating Request: %v", err.Error())
		}
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)
		assert.Contains(t, res.Body.String(), "Internal Server Error")
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})
	mockedDb.EXPECT().ResetPassword("12340987wertt", gomock.Any()).Return(nil)
	t.Run("Test_For_Reset_Password_Error", func(t *testing.T) {
		body, _ := json.Marshal(requestPassword)
		req, err := http.NewRequest(http.MethodPost, "/api/v1/reset-password/12340987wertt", strings.NewReader(string(body)))
		if err != nil {
			t.Errorf("Error Creating Request: %v", err.Error())
		}
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)
		assert.Contains(t, res.Body.String(), "Password Reset Successfully")
		assert.Equal(t, http.StatusOK, res.Code)
	})
}
