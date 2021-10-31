package server

import (
	"github.com/decadevs/rentals-api/models"
	"github.com/decadevs/rentals-api/server/response"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (s *Server) handleGetAllApartments() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get all apartments should only be accessible by those that have the permission
		if userI, exists := c.Get("user"); exists {
			if user, ok := userI.(*models.User); ok {
				//userId := c.Param("userId")
				apartment, err := s.DB.GetAllApartments(user.ID)
				if err != nil {
					log.Printf("get apartments error : %v\n", err)
					response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
					return
				}
				response.JSON(c, "retrieved apartments successfully", http.StatusOK, gin.H{"apartment": apartment}, nil)
				return
			}
		}
		log.Printf("can't get apartment from context\n")
		response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
		return
	}
}
