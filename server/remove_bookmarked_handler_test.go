package server

import (
	"fmt"
	"github.com/decadevs/rentals-api/db"
	"github.com/decadevs/rentals-api/models"
	"github.com/decadevs/rentals-api/router"
	"github.com/decadevs/rentals-api/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestApplication_RemoveBookmarked(t *testing.T) {
	ctr := gomock.NewController(t)
	mockedDB := db.NewMockDB(ctr)

	s := &Server{
		DB:     mockedDB,
		Router: router.NewRouter(),
	}
	router := s.setupRouter()
	user := &models.User{
		Email: "shuaib@gmail.com",
	}
	user.ID = "1234567asdf"
	secret := os.Getenv("JWT_SECRET")
	accessClaims, refreshClaimns := services.GenerateClaims(user.Email)
	accToken, _ := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
	services.GenerateToken(jwt.SigningMethodHS256, refreshClaimns, &secret)

	url := "/api/v1/user/156uhjqhacgyqfa/removebookmark"

	mockedDB.EXPECT().TokenInBlacklist(accToken).Return(false).Times(3)
	mockedDB.EXPECT().FindUserByEmail(user.Email).Return(user, nil).Times(3)

	t.Run("Test_For_Already_Bookmarked_Apartment", func(t *testing.T) {

		mockedDB.EXPECT().CheckApartmentInBookmarkApartment("1234567asdf", "156uhjqhacgyqfa").Return(false)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))

		router.ServeHTTP(w, req)

		log.Println(w.Body)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "apartment not bookmarked")
	})

	t.Run("Test_For_Error_Removing_BookmarkedApartment", func(t *testing.T) {
		apartment := &models.BookmarkApartment{
			UserID:      "1234567asdf",
			ApartmentID: "156uhjqhacgyqfa",
		}

		mockedDB.EXPECT().CheckApartmentInBookmarkApartment("1234567asdf", "156uhjqhacgyqfa").Return(true)
		mockedDB.EXPECT().RemoveBookmarkedApartment(apartment).Return(errors.New("an error occurred"))
		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))

		router.ServeHTTP(w, req)

		log.Println(w.Body)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "internal server error")
	})

	t.Run("Test_For_Success", func(t *testing.T) {
		apartment := &models.BookmarkApartment{
			UserID:      "1234567asdf",
			ApartmentID: "156uhjqhacgyqfa",
		}

		mockedDB.EXPECT().CheckApartmentInBookmarkApartment("1234567asdf", "156uhjqhacgyqfa").Return(true)
		mockedDB.EXPECT().RemoveBookmarkedApartment(apartment).Return(nil)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))

		router.ServeHTTP(w, req)

		log.Println(w.Body)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Bookmarked Remove Successfully")
	})
}