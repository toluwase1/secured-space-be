package server

import (
	"fmt"
	"github.com/decadevs/rentals-api/server/response"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (s *Server) VerifyEmail() gin.HandlerFunc {
	return func(c *gin.Context) {

		ID := c.Param("userID")
		token := c.Param("userToken")

		_ , err := s.DB.FindUserByID(ID)
		if err != nil{
			log.Printf("Error: %v", err.Error())
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"User not  found"})
			return
		}

		user, err := s.DB.CompareToken(ID)
		if err != nil{
			log.Printf("Error: %v", err.Error())
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"Internal server error"})
			return
		}

		if token != user.Token{
			log.Println("invalid token")
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"Invalid user token or ID"})
			return
		}

		err = s.DB.SetUserToActive(ID)
		if err != nil{
			log.Printf("Error: %v", err.Error())
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"Could not set user"})
			return
		}
		response.JSON(c, fmt.Sprintf("%s,your email has been verified successfully.",user.FirstName), http.StatusOK, nil, nil)
	}
}