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
			return
		}
		userId := userI.(*models.User).ID
		apartment.UserID = userId

		if err := s.decode(c, &apartment); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}

		err := s.DB.CreateApartment(&apartment)
		if err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, []string{err.Error()})
			return
		}
		response.JSON(c, "Apartment Successfully Added", http.StatusOK, apartment, nil)

	}
}

func (s *Server) DeleteApartment() gin.HandlerFunc {
	return func(c *gin.Context) {
		if userI, exists := c.Get("user"); exists {
			if user, ok := userI.(*models.User); ok {
				apartmentID := c.Param("apartmentID")
				if apartmentID == "" {
					response.JSON(c, "", http.StatusBadRequest, nil, []string{"apartment id cannot be empty"})
					return
				}
				err := s.DB.DeleteApartment(apartmentID, user.ID)
				if err != nil {
					log.Println(err.Error())
					response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
					return
				}
				response.JSON(c, "Deleted Successfully", http.StatusOK, nil, nil)
				return
			}
		}
		log.Printf("can't get user from context\n")
		response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
	}
}
func (s *Server) handleUpdateApartmentDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		apartmentID := c.Param("apartmentID")
		if apartmentID == "" {
			response.JSON(c, "", http.StatusBadRequest, nil, []string{"apartment id cannot be empty"})
			return
		}
		apartment := &models.Apartment{}
		if errs := s.decode(c, apartment); errs != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, errs)
			return
		}

		if err := s.DB.UpdateApartment(apartment, apartmentID); err != nil {
			log.Printf("update apartment error : %v\n", err)
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
			return
		}
		response.JSON(c, "apartment updated successfully", http.StatusOK, nil, nil)
		return
	}
}
