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
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func AuthorizeTestRoutes(m *db.MockDB, t *testing.T) *string {
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
func TestServer_SaveBookmarkApartment(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := db.NewMockDB(ctrl)

	s := &Server{
		DB:     m,
		Router: r.NewRouter(),
	}

	router := s.setupRouter()

	t.Run("Test_For_Empty_ApartmentID", func(t *testing.T) {

		accToken := AuthorizeTestRoutes(m, t)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/user//bookmark", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))

		router.ServeHTTP(w, req)

		log.Println(w.Body)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "apartment id cannot be empty")
	})

	t.Run("Test_For_Already_Existing_Apartment", func(t *testing.T) {

		accToken := AuthorizeTestRoutes(m, t)
		m.EXPECT().CheckApartmentInBookmarkApartment("1234567asdf", "14uhjqhacgyqfa").Return(true)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/user/14uhjqhacgyqfa/bookmark", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))

		router.ServeHTTP(w, req)

		log.Println(w.Body)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "already bookmarked this apartment")
	})

	t.Run("Test_For_Error_Saving_BookmarkApartment", func(t *testing.T) {
		apartment := &models.BookmarkApartment{
			UserID:      "1234567asdf",
			ApartmentID: "14uhjqhacgyqfa",
		}
		accToken := AuthorizeTestRoutes(m, t)
		m.EXPECT().CheckApartmentInBookmarkApartment("1234567asdf", "14uhjqhacgyqfa").Return(false)
		m.EXPECT().SaveBookmarkApartment(apartment).Return(errors.New("error"))
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/user/14uhjqhacgyqfa/bookmark", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))

		router.ServeHTTP(w, req)

		log.Println(w.Body)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "internal server error")
	})

	t.Run("Test_For_Success", func(t *testing.T) {
		apartment := &models.BookmarkApartment{
			UserID:      "1234567asdf",
			ApartmentID: "14uhjqhacgyqfa",
		}
		accToken := AuthorizeTestRoutes(m, t)
		m.EXPECT().CheckApartmentInBookmarkApartment("1234567asdf", "14uhjqhacgyqfa").Return(false)
		m.EXPECT().SaveBookmarkApartment(apartment).Return(nil)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/user/14uhjqhacgyqfa/bookmark", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))

		router.ServeHTTP(w, req)

		log.Println(w.Body)
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "Saved Successfully")
	})
}
