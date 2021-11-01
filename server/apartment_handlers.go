package server

import (
	"github.com/decadevs/rentals-api/models"
	"github.com/decadevs/rentals-api/server/response"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

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
