package server

import (
	"fmt"
	"github.com/decadevs/rentals-api/server/response"
	"github.com/decadevs/rentals-api/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (s *Server) ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		resetPassword := struct {
			Password string `json:"password" binding:"required"`
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

func (s *Server) ForgotPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := struct{
			Email string `json:"email" binding:"required"`
		}{}
		errs := s.decode(c, &email)
		if errs != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, errs)
            return
		}
		user, err := s.DB.FindUserByEmail(email.Email)
		if err != nil {
			log.Printf("Error: %v", err.Error())
            response.JSON(c, "", http.StatusBadRequest, nil, []string{"email does not exist"})
            return
		}
		_, err = s.Mail.SendResetPassword(user.Email, fmt.Sprintf("http://localhost:3000/reset-password/%s", user.ID))
		if err != nil {
			log.Printf("Error: %v", err.Error())
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"Internal Server Error"})
			return
		}
		response.JSON(c, "Reset Password Link Sent Successfully", http.StatusOK, nil, nil)
	}
}
