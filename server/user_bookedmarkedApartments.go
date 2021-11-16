package server

import (
	"github.com/decadevs/rentals-api/server/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) GetBookmarkedApartments() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUserFromContext(c)
		if err != nil {
			response.JSON(c, "", http.StatusUnauthorized, nil, []string{"user not authorized"})
			return
		}
		bookmarkedApartment, err := s.DB.GetBookmarkedApartments(user.ID)
		if err != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
			return
		}
		response.JSON(c, "retrieved bookmarks successfully", http.StatusOK, bookmarkedApartment, nil)
		return
	}
}
