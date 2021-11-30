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
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestServer_DeleteApartment(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println(err.Error())
	}
	ctrl := gomock.NewController(t)
	mockDB := db.NewMockDB(ctrl)
	s := &Server{
		DB:     mockDB,
		Router: router.NewRouter(),
	}
	r := s.setupRouter()

	accessClaims, _ := services.GenerateClaims("adebayo@gmail.com")

	secret := os.Getenv("JWT_SECRET")
	accToken, err := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
	if err != nil {
		t.Fail()
	}
	user := &models.User{Email: "adebayo@gmail.com"}
	user.ID = "aefrfh123435waes"
	mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false).Times(2)
	mockDB.EXPECT().FindUserByEmail(user.Email).Return(user, nil).Times(2)
	t.Run("Test_Empty_ApartmentID", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/api/v1/user/apartment/", nil)
		if err != nil {
			t.Fail()
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusNotFound, rw.Code)
		assert.Contains(t, rw.Body.String(), "404 page not found")
	})
	mockDB.EXPECT().DeleteApartment("12323shjbvbhj1t", user.ID).Return(errors.New("an error occurred"))
	t.Run("Test_Error-Deleting_Apartment", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/api/v1/user/apartment/12323shjbvbhj1t", nil)
		if err != nil {
			t.Fail()
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		assert.Contains(t, rw.Body.String(), "internal server error")
	})

	mockDB.EXPECT().DeleteApartment("12323shjbvbhj1t", user.ID).Return(nil)
	t.Run("Test_For_Successful_Delete", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/api/v1/user/apartment/12323shjbvbhj1t", nil)
		if err != nil {
			t.Fail()
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), "Deleted Successfully")
	})
}

func TestApplication_GetApartmentDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	mdb := db.NewMockDB(ctrl)
	apartmentID := "12sdfg-456hcvbn-ut78okjh"

	s := &Server{
		DB:     mdb,
		Router: router.NewRouter(),
	}

	route := s.setupRouter()
	apartment := &models.Apartment{
		UserID:          "dddfehrrbb1t32447",
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

	mdb.EXPECT().ApartmentDetails(apartmentID).Return(nil, errors.New("error exist"))
	mdb.EXPECT().ApartmentDetails(apartmentID).Return(nil, nil)

	t.Run("testing error", func(t *testing.T) {

		jsonapartment, err := json.Marshal(apartment)
		if err != nil {
			t.Fail()
			return
		}
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/apartment-details/%s", apartmentID), strings.NewReader(string(jsonapartment)))
		req.Header.Set("Content-Type", "application/json")
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
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/apartment-details/%s", apartmentID), strings.NewReader(string(jsonapartment)))
		req.Header.Set("Content-Type", "application/json")
		route.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "apartment retrieved successfully")

	})
}
