package server

import (
	"github.com/decadevs/rentals-api/models"
	"github.com/decadevs/rentals-api/server/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) handleCreateApartment() gin.HandlerFunc {
	// function to handle adding an apartment
	return func(c *gin.Context) {
		apartment := models.Apartment{}

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
