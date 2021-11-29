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

//func Test_CreateApartment(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	mockedDB := db.NewMockDB(ctrl)
//
//	godotenv.Load("../.env")
//	s := &Server{
//		mockedDB,
//		router.NewRouter(),
//	}
//	route := s.setupRouter()
//	user := &models.User{Email: "boboti@gmail.com"}
//	apartment := &models.Apartment{
//		Models:          models.Models{},
//		UserID:          "esdfghjkfhgjh",
//		Title:           "Fantastic family flat",
//		CategoryID:      "qesrdfhgjhkj",
//		Description:     "3 bedroom with awesome neighbors",
//		Price:           40000,
//		NoOfRooms:       3,
//		Furnished:       false,
//		Location:        "Lagos",
//		ApartmentStatus: false,
//		Images:          nil,
//		Interiors:       nil,
//		Exteriors:       nil,
//	}
//	marshalledApartment, _ := json.Marshal(apartment)
//
//	secret := os.Getenv("JWT_SECRET")
//	accessClaims, refreshClaims := services.GenerateClaims(user.Email)
//	accToken, _ := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
//	services.GenerateToken(&jwt.SigningMethodHMAC{}, refreshClaims, &secret)
//
//	mockedDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false).Times(2)
//	mockedDB.EXPECT().FindUserByEmail(user.Email).Return(user, nil).Times(2)
//
//	mockedDB.EXPECT().CreateApartment(apartment).Return(nil)
//	t.Run("Testing_For_Apartment_Successfully_Added", func(t *testing.T) {
//		rw := httptest.NewRecorder()
//		req, _ := http.NewRequest(http.MethodPost, "/api/v1/user/apartments", strings.NewReader(string(marshalledApartment)))
//		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
//		route.ServeHTTP(rw, req)
//
//		assert.Equal(t, http.StatusOK, rw.Code)
//		assert.Contains(t, rw.Body.String(), "Apartment Successfully Added")
//	})
//
//	mockedDB.EXPECT().CreateApartment(apartment).Return(errors.New("error creating apartment"))
//	t.Run("Testing_For_Error_in_Creating_Apartment", func(t *testing.T) {
//		rw := httptest.NewRecorder()
//		req, _ := http.NewRequest(http.MethodPost, "/api/v1/user/apartments", strings.NewReader(string(marshalledApartment)))
//		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
//		route.ServeHTTP(rw, req)
//
//		assert.Equal(t, http.StatusBadRequest, rw.Code)
//		assert.Contains(t, rw.Body.String(), "Bad Request")
//	})
//}

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
	mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false).Times(3)
	mockDB.EXPECT().FindUserByEmail(user.Email).Return(user, nil).Times(3)
	t.Run("Test_Empty_ApartmentID", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/api/v1/user/apartment//", nil)
		if err != nil {
			t.Fail()
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		assert.Contains(t, rw.Body.String(), "apartment id cannot be empty")
	})
	mockDB.EXPECT().DeleteApartment("12323shjbvbhj1t", user.ID).Return(errors.New("an error occurred"))
	t.Run("Test_Error-Deleting_Apartment", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/api/v1/user/apartment/12323shjbvbhj1t/", nil)
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
		req, err := http.NewRequest(http.MethodDelete, "/api/v1/user/apartment/12323shjbvbhj1t/", nil)
		if err != nil {
			t.Fail()
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), "Deleted Successfully")
	})
}

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
	accessClaims, _ := services.GenerateClaims("franklyn@gmail.com")
	accToken, _ := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
	log.Println(*accToken)
	route := s.setupRouter()
	apartment := &models.Apartment{
		User: models.User{
			FirstName: "a",
			LastName:  "b",
			Email:     "c@gmail.com",
			Phone1:    "123",
			Password:  "1223",
		},
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
	m.EXPECT().FindUserByEmail("franklyn@gmail.com").Return(user, nil).Times(3)
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
