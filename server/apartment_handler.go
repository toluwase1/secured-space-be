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

		// get the user id from a logged-in user
		user, exists:= c.Get("user")
		if !exists{
			log.Printf("can't get user from context\n")
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
		}
		userId := user.(*models.User).ID
		apartment.UserID = userId

		if errs := s.decode(c, apartment); errs != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, errs)
			return
		}
		_, err := s.DB.CreateApartment(&apartment)
		if err != nil{
			response.JSON(c,"",http.StatusBadRequest,nil,[]string{err.Error()})
		}
		response.JSON(c,"Apartment Successfully Added",http.StatusOK,nil,nil)

	}
}
