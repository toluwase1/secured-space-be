package server

import (
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
				log.Printf("Error: %v", err.Error())
				response.JSON(c, "", http.StatusInternalServerError, nil, []string{"Invalid user token"})
				return
		}

		err = s.DB.SetUserToActive(ID)
		if err != nil{
			log.Printf("Error: %v", err.Error())
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"Could not set user"})
			return
		}
		response.JSON(c, "user verified successfully", http.StatusOK, nil, nil)
	}
}