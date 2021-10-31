package server

import (
	"github.com/decadevs/rentals-api/models"
	"github.com/decadevs/rentals-api/server/response"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (s *Server) handleCreateApartment() gin.HandlerFunc {
	// function to handle adding an apartment
	return func(c *gin.Context) {
		apartment := models.Apartment{}

		//get the user id from a logged-in user
		userI, exists := c.Get("user")
		if !exists {
			log.Printf("can't get user from context\n")
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
		}
		userId := userI.(*models.User).ID
		apartment.UserID = userId

		if err := s.decode(c, &apartment); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}

		_, err := s.DB.CreateApartment(&apartment)
		if err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, []string{err.Error()})
		}
		response.JSON(c, "Apartment Successfully Added", http.StatusOK, nil, nil)

	}
}
