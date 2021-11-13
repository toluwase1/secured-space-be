package server

import (
	"github.com/decadevs/rentals-api/server/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUserFromContext(c)
		if err != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"Internal Server Error"})
		}
	}
}