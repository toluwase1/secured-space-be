package server

import (
	"encoding/json"
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

	s := &Server{
		DB:     m,
		Router: router.NewRouter(),
	}
	router := s.setupRouter()
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

	m.EXPECT().UpdateApartment(apartment).Return(nil)
	t.Run("testing it no error", func(t *testing.T) {
		jsonapartment, err := json.Marshal(apartment)
		if err != nil {
			t.Fail()
			return
		}
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/user/apartment/update", strings.NewReader(string(jsonapartment)))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "apartment updated successfully")
	})

	m.EXPECT().UpdateApartment(apartment).Return(errors.New("error exist"))
	t.Run("testing error", func(t *testing.T) {
		jsonapartment, err := json.Marshal(apartment)
		if err != nil {
			t.Fail()
			return
		}
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/user/apartment/update", strings.NewReader(string(jsonapartment)))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "internal server error")
	})

}
