package server

import (
	"encoding/json"
	"github.com/decadevs/rentals-api/db"
	"github.com/decadevs/rentals-api/models"
	"github.com/decadevs/rentals-api/services"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_handleLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedDB := db.NewMockDB(ctrl)

	s := &Server{
		DB: mockedDB,
	}

	router := s.setupRouter()

	t.Run("Test_For_Login_Request", func(t *testing.T) {
		loginRequest := &struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}{
			Email:    "jdoe@gmail.com",
			Password: "",
		}
		jsonFile, err := json.Marshal(loginRequest)
		if err != nil {
			t.Fail()
		}
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(string(jsonFile)))

		router.ServeHTTP(w, req)

		log.Println(w.Body.String())
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "validation failed on field 'Password', condition: required")
	})


	t.Run("Test_FindUserByEmail", func(t *testing.T) {
		hashedP, _ := services.GenerateHashPassword("password")
		user := &models.User{Email: "jdoe@gmail.com", HashedPassword: string(hashedP)}
		mockedDB.EXPECT().FindUserByEmail(user.Email).Return(user, nil)
		loginRequest := &struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}{
			Email:    "jdoe@gmail.com",
			Password: "password",
		}
		jsonFile, err := json.Marshal(loginRequest)
		if err != nil {
			t.Fail()
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(string(jsonFile)))


		router.ServeHTTP(w, req)
		log.Println(w.Body)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "")
	})
}

//func TestAuthoriz(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	m := db.NewMockDB(ctrl)
//
//	s := &Server{
//		DB:     m,
//		Router: router.NewRouter(),
//	}
//	router := s.setupRouter()
//
//	t.Run("Test_Authorize", func(t *testing.T) {
//		accessClaims, _ := services.GenerateClaims("adebayo@gmail.com")
//
//		secret := os.Getenv("JWT_SECRET")
//		accToken, err := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
//		if err != nil {
//			t.Fail()
//		}
//
//		user := &models.User{Email: "adebayo@gmail.com"}
//		users := []models.User{{}}
//		m.EXPECT().FindUserByEmail(user.Email).Return(user, nil)
//		m.EXPECT().TokenInBlacklist(accToken).Return(false)
//		m.EXPECT().FindAllUsersExcept(user.Email).Return(users, nil)
//
//		w := httptest.NewRecorder()
//		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
//		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
//
//		router.ServeHTTP(w, req)
//
//		assert.Equal(t, http.StatusOK, w.Code)
//		assert.Contains(t, w.Body.String(), "retrieved users successfully")
//	})
//
//	t.Run("Test_FIndUserByEmail_Error", func(t *testing.T) {
//		accessClaims, _ := services.GenerateClaims("adebayo@gmail.com")
//
//		secret := os.Getenv("JWT_SECRET")
//		accToken, err := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
//		if err != nil {
//			t.Fail()
//		}
//
//		user := &models.User{Email: "adebayo@gmail.com"}
//		m.EXPECT().FindUserByEmail(user.Email).Return(user, errors.New("an error occurred"))
//		m.EXPECT().TokenInBlacklist(accToken).Return(false)
//
//		w := httptest.NewRecorder()
//		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
//		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
//
//		router.ServeHTTP(w, req)
//		assert.Equal(t, http.StatusNotFound, w.Code)
//		assert.Contains(t, w.Body.String(), "user not found")
//	})
//}

