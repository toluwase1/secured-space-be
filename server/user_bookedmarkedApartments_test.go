package server

import (
	"fmt"
	"github.com/decadevs/rentals-api/db"
	"github.com/decadevs/rentals-api/models"
	r "github.com/decadevs/rentals-api/router"
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

func TestApplication_BookmarkedApartments(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedDB := db.NewMockDB(ctrl)

	s := &Server{
		DB:     mockedDB,
		Router: r.NewRouter(),
	}
	router := s.setupRouter()
	user := &models.User{
		Email: "adebayo@gmail.com",
	}
	//bookmarks := &models.BookmarkApartment{
	//	UserID: "techagentng@gmail.com",
	//	ApartmentID: "q2342342ccwefcwef",
	//}
	user.ID = "cad3ac6e-ad3d-46c4-8b20-4283525c6136"
	secret := os.Getenv("JWT_SECRET")
	accessClaims, refreshClaims := services.GenerateClaims(user.Email)
	accToken, _ := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
	services.GenerateToken(jwt.SigningMethodHS256, refreshClaims, &secret)

	mockedDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false).Times(2)
	mockedDB.EXPECT().FindUserByEmail(user.Email).Return(user, nil).Times(2)
	//mockedDB.EXPECT().GetBookmarkedApartment(user.ID).Return(user.ID, nil).Times(2)

	mockedDB.EXPECT().GetBookmarkedApartments(user.ID).Return(nil, nil)
	t.Run("Test_For_Successful_Response", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/bookmark/apartments", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		router.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), "retrieved bookmarks successfully")
	})

	mockedDB.EXPECT().GetBookmarkedApartments(user.ID).Return(nil, errors.New("An error occurred"))

	t.Run("Test_For_Error_in_Getting_Bookmarked_Apartment", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/bookmark/apartments", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		router.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		assert.Contains(t, rw.Body.String(), "internal server error")
	})
}