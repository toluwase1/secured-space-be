package server

import (
	"github.com/decadevs/rentals-api/server/response"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (s *Server) SearchApartment() gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryID := c.Query("category_id")
		location := c.Query("location")
		minPrice := c.Query("min_price")
		maxPrice := c.Query("max_price")
		noOfRooms := c.Query("no_of_rooms")

	    apartments, err := s.DB.SearchApartment(categoryID, location, minPrice, maxPrice, noOfRooms)
		if err != nil {
			log.Printf("Error: %v", err.Error())
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"Internal Server Error"})
			return
		}
		response.JSON(c, "", http.StatusOK, apartments, nil)
	}
}
