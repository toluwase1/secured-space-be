package server

import (
	"github.com/decadevs/rentals-api/server/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		resetPassword:=  struct{
			Password	string	`json:"password" binding:"required"`
			ConfirmPassword	string	`json:"confirm_password" binding:"required"`
		}{}
		errs := s.decode(c, &resetPassword)
		if errs != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, errs)
			return
		}
		userID := c.Param("userID")
		if userID == "" {
            response.JSON(c, "", http.StatusBadRequest, nil, []string{"user id is required"})
            return
        }
		if resetPassword.Password != resetPassword.ConfirmPassword {
			response.JSON(c, "", http.StatusBadRequest, nil, []string{"Password do not match"})
			return
		}

	}
}
