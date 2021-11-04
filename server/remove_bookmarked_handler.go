package server

import (
	"github.com/decadevs/rentals-api/models"
	"github.com/decadevs/rentals-api/server/response"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func (s *Server) RemoveBookmarkedApartment() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, exists := c.Get("user"); exists {
			if user, ok := user.(*models.User); ok {
				apartmentID := strings.TrimSpace(c.Param("apartmentID"))

				if ok := s.DB.CheckApartmentInBookmarkApartment(user.ID, apartmentID); !ok {
					response.JSON(c, "", http.StatusBadRequest, nil, []string{"apartment not bookmarked"})
					return
				}
				rmbookmarkApartment := &models.BookmarkApartment{
					UserID:      user.ID,
					ApartmentID: apartmentID,
				}
				if err := s.DB.RemoveBookmarkedApartment(rmbookmarkApartment); err != nil {
					log.Printf(err.Error())
					response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
					return
				}
				response.JSON(c, "Bookmarked Remove  Successfully", http.StatusOK, nil, nil)
				return
			}
		}
		log.Printf("can't get user from context\n")
		response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
	}
}