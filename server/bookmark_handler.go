package server

import (
	"github.com/decadevs/rentals-api/models"
	"github.com/decadevs/rentals-api/server/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) SaveBookmarkApartment() gin.HandlerFunc {
	return func(c *gin.Context) {
		bookmarkApartment := &models.BookmarkApartment{}

		if errs := s.decode(c, bookmarkApartment); errs != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, errs)
			return
		}

	}
}
