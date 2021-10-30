package server

import (
	"encoding/json"
	"fmt"
	"github.com/decadevs/rentals-api/db"
	"github.com/decadevs/rentals-api/models"
	"github.com/decadevs/rentals-api/router"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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

	m.EXPECT().UpdateApartment(apartment, apartmentID).Return(nil)
	t.Run("testing empty apartment id", func(t *testing.T) {

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/user/%s/update", ""), nil)
		route.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "apartment id cannot be empty")
	})

	t.Run("testing if error does not exist", func(t *testing.T) {
		jsonapartment, err := json.Marshal(apartment)
		if err != nil {
			t.Fail()
			return
		}
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/user/%s/update", apartmentID), strings.NewReader(string(jsonapartment)))
		route.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "apartment updated successfully")
	})

	m.EXPECT().UpdateApartment(apartment, apartmentID).Return(errors.New("error exist"))
	t.Run("testing error", func(t *testing.T) {
		jsonapartment, err := json.Marshal(apartment)
		if err != nil {
			t.Fail()
			return
		}
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/user/%s/update", apartmentID), strings.NewReader(string(jsonapartment)))
		route.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "internal server error")
	})

}
