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
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func CheckAuthorize(m *db.MockDB, t *testing.T) *string {
	accessClaims, _ := services.GenerateClaims("adebayo@gmail.com")

	secret := os.Getenv("JWT_SECRET")
	accToken, err := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
	if err != nil {
		t.Fail()
	}

	// "password" was hashed to form the hashedPassword data, and it is for testing purpose
	user := &models.User{Email: "adebayo@gmail.com", HashedPassword: "$2a$10$dP0JsmQr4ycXj8MQJVaDkedmOA2owa7lAOKUqOmwEN3IbDOfkLROy"}
	user.ID = "1234567asdf"
	m.EXPECT().FindUserByEmail(user.Email).Return(user, nil)
	m.EXPECT().TokenInBlacklist(accToken).Return(false)
	return accToken
}
func TestApplication_ChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedDb := db.NewMockDB(ctrl)

	s := &Server{
		DB:     mockedDb,
		Router: router.NewRouter(),
	}
	r := s.setupRouter()
	accToken := CheckAuthorize(mockedDb, t)
	t.Run("Test_For_Missing_Password", func(t *testing.T) {
		passwordInfo := struct {
			OldPassword string `json:"old_password" binding:"required"`
			NewPassword string `json:"new_password" binding:"required"`
		}{
			OldPassword: "",
			NewPassword: "",
		}
		body, _ := json.Marshal(passwordInfo)
		req, err := http.NewRequest(http.MethodPost, "/api/v1/user/change-password", strings.NewReader(string(body)))
		if err != nil {
			t.Errorf("Error Creating Request: %v", err.Error())
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)
		assert.Contains(t, res.Body.String(), "validation failed on field 'OldPassword', condition: required")
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})
	accToken = CheckAuthorize(mockedDb, t)
	t.Run("Test_For_Incorrect_Password", func(t *testing.T) {
		passwordInfo := struct {
			OldPassword string `json:"old_password" binding:"required"`
			NewPassword string `json:"new_password" binding:"required"`
		}{
			OldPassword: "password123",
			NewPassword: "password",
		}
		body, _ := json.Marshal(passwordInfo)
		req, err := http.NewRequest(http.MethodPost, "/api/v1/user/change-password", strings.NewReader(string(body)))
		if err != nil {
			t.Errorf("Error Creating Request: %v", err.Error())
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)
		assert.Contains(t, res.Body.String(), "Incorrect Password Details")
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	accToken = CheckAuthorize(mockedDb, t)
	mockedDb.EXPECT().ResetPassword("1234567asdf", gomock.Any()).Return(errors.New("an error occurred"))
	t.Run("Test_For_Password_Changed_Error", func(t *testing.T) {
		passwordInfo := struct {
			OldPassword string `json:"old_password" binding:"required"`
			NewPassword string `json:"new_password" binding:"required"`
		}{
			OldPassword: "password",
			NewPassword: "password123",
		}
		body, _ := json.Marshal(passwordInfo)
		req, err := http.NewRequest(http.MethodPost, "/api/v1/user/change-password", strings.NewReader(string(body)))
		if err != nil {
			t.Errorf("Error Creating Request: %v", err.Error())
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)
		assert.Contains(t, res.Body.String(), "Internal Server Error")
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})
	accToken = CheckAuthorize(mockedDb, t)
	mockedDb.EXPECT().ResetPassword("1234567asdf", gomock.Any()).Return(nil)
	t.Run("Test_For_Successfully_Password_Change", func(t *testing.T) {
		passwordInfo := struct {
			OldPassword string `json:"old_password" binding:"required"`
			NewPassword string `json:"new_password" binding:"required"`
		}{
			OldPassword: "password",
			NewPassword: "password123",
		}
		body, _ := json.Marshal(passwordInfo)
		req, err := http.NewRequest(http.MethodPost, "/api/v1/user/change-password", strings.NewReader(string(body)))
		if err != nil {
			t.Errorf("Error Creating Request: %v", err.Error())
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)
		assert.Contains(t, res.Body.String(), "Password Changed Successfully")
		assert.Equal(t, http.StatusOK, res.Code)
	})
}
