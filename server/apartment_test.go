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

func TestUpdateApartment(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := db.NewMockDB(ctrl)
	apartmentID := "12sdfg-456hcvbn-ut78okjh"

	s := &Server{
		DB:     m,
		Router: router.NewRouter(),
	}
	user := &models.User{}
	secret := os.Getenv("JWT_SECRET")
	accessClaims, _ := services.GenerateClaims("frank@gmail.com")
	accToken, _ := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)

	route := s.setupRouter()
	apartment := &models.Apartment{
		Title:           "2 bedrooms",
		Description:     "Bay area lodge",
		Price:           45000,
		NoOfRooms:       3,
		Furnished:       false,
		Location:        "lagos",
		ApartmentStatus: false,
		Interiors:       []models.InteriorFeature{{Name: "gym"}, {Name: "fire place"}},
		Exteriors:       []models.ExteriorFeature{{Name: "Garage"}, {Name: "pool"}},
	}

	m.EXPECT().TokenInBlacklist(gomock.Any()).Return(false).Times(3)
	m.EXPECT().FindUserByEmail("frank@gmail.com").Return(user, nil).Times(3)
	m.EXPECT().UpdateApartment(apartment, apartmentID).Return(errors.New("error exist"))
	m.EXPECT().UpdateApartment(apartment, apartmentID).Return(nil)

	t.Run("testing empty apartment id", func(t *testing.T) {

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/user/%s/update", ""), nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		route.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "apartment id cannot be empty")
	})

	t.Run("testing error", func(t *testing.T) {

		jsonapartment, err := json.Marshal(apartment)
		if err != nil {
			t.Fail()
			return
		}
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/user/%s/update", apartmentID), strings.NewReader(string(jsonapartment)))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		route.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "internal server error")
	})

	t.Run("testing if error does not exist", func(t *testing.T) {
		jsonapartment, err := json.Marshal(apartment)
		if err != nil {
			t.Fail()
			return
		}
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/user/%s/update", apartmentID), strings.NewReader(string(jsonapartment)))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		route.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "apartment updated successfully")
	})
}
