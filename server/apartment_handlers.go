package server

import (
	"github.com/decadevs/rentals-api/server/response"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (s *Server) handleUpdateApartmentDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Get("apartment")

		log.Printf("can't get apartment from context\n")
		response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
	}
}
