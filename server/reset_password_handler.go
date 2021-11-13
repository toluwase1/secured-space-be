package server

import (
	"github.com/decadevs/rentals-api/server/response"
	"github.com/decadevs/rentals-api/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (s *Server) ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		resetPassword:=  struct{
			Password	string	`json:"password" binding:"required"`
		}{}
		errs := s.decode(c, &resetPassword)
		if errs != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, errs)
			return
		}
		userID := c.Param("userID")
		hashedPassword, err := services.GenerateHashPassword(resetPassword.Password)
		if err != nil {
			log.Printf("Error: %v", err.Error())
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"Internal Server Error"})
			return
		}
		err = s.DB.ResetPassword(userID, string(hashedPassword))
		if err != nil {
			log.Printf("Error: %v", err.Error())
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"Internal Server Error"})
			return
		}
		response.JSON(c, "Password Reset Successfully", http.StatusOK, nil, nil)
	}
}
