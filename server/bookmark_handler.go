package server

import (
	"github.com/decadevs/rentals-api/models"
	"github.com/decadevs/rentals-api/server/response"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func (s *Server) SaveBookmarkApartment() gin.HandlerFunc {
	return func(c *gin.Context) {
		if userI, exists := c.Get("user"); exists {
			if user, ok := userI.(*models.User); ok {
				apartmentID := strings.TrimSpace(c.Param("apartmentID"))
				if apartmentID == "" {
					response.JSON(c, "", http.StatusBadRequest, nil, []string{"apartment id cannot be empty"})
					return
				}
				if ok := s.DB.CheckApartmentInBookmarkApartment(user.ID, apartmentID); ok {
					response.JSON(c, "", http.StatusBadRequest, nil, []string{"already bookmarked this apartment"})
					return
				}
				bookmarkApartment := &models.BookmarkApartment{
					UserID:      user.ID,
					ApartmentID: apartmentID,
				}
				if err := s.DB.SaveBookmarkApartment(bookmarkApartment); err != nil {
					log.Printf(err.Error())
					response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
					return
				}
				response.JSON(c, "Saved Successfully", http.StatusCreated, nil, nil)
				return
			}
		}
		log.Printf("can't get user from context\n")
		response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
	}
}
