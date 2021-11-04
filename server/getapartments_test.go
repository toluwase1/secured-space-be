package server

import (
	"fmt"
	"github.com/decadevs/rentals-api/db"
	"github.com/decadevs/rentals-api/router"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)



func TestApplication_GetUserApartments(t *testing.T) {
	ctrl := gomock.NewController(t)
	mDB := db.NewMockDB(ctrl)

	s := &Server{
		DB:     mDB,
		Router: router.NewRouter(),
	}
	router := s.setupRouter()

	t.Run("Testing_Error_Getting_All_Apartments", func(t *testing.T) {
		user, accToken := AuthorizeTestRoute(mDB, t, "tchisom@gmail.com")
		mDB.EXPECT().GetUsersApartments(user.ID).Return(nil, errors.New("an error occurred"))
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/user-apartment", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		router.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		assert.Contains(t, rw.Body.String(), "internal server error")
	})


	t.Run("Test_For_Successful_Response", func(t *testing.T) {
		user, accToken := AuthorizeTestRoute(mDB, t, "tchisom@gmail.com")
		mDB.EXPECT().GetUsersApartments(user.ID).Return(nil, nil)
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/user-apartment", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		router.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), "retrieved apartments successfully")
	})
}


