package server

import (
	"github.com/decadevs/rentals-api/models"
	"github.com/decadevs/rentals-api/server/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetBookmarkedApartments lists all bookmarked apartments
func (s *Server) GetBookmarkedApartments() gin.HandlerFunc {
	return func(c *gin.Context) {
		if userI, exists := c.Get("user"); exists {
			if user, ok := userI.(*models.User); ok {
				bookmarkedApartment, err := s.DB.GetBookmarkedApartments(user.ID)
				if err != nil {
					response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
					return
				}
				response.JSON(c, "retrieved bookmarks successfully", http.StatusOK, bookmarkedApartment, nil)
			}
		}
	}
}