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
		Router: router.NewRouter(),
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
