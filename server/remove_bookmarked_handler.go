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
		user, err := GetUserFromContext(c)

		if err != nil {
			response.JSON(c, "", http.StatusUnauthorized, nil, []string{"user not authorized"})
		}
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

