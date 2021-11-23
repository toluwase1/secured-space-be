package server

import (
	"github.com/decadevs/rentals-api/server/response"
	"github.com/decadevs/rentals-api/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUserFromContext(c)
		if err != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"Internal Server Error"})
			return
		}
		passwordInfo := struct {
			OldPassword string `json:"old_password" binding:"required"`
			NewPassword string `json:"new_password" binding:"required"`
		}{}
		if errs := s.decode(c, &passwordInfo); errs != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, errs)
			return
		}
		err = services.CompareHashAndPassword([]byte(passwordInfo.OldPassword), user.HashedPassword)
		if err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, []string{"Incorrect Password Details"})
			return
		}
		hashedPassword, err := services.GenerateHashPassword(passwordInfo.NewPassword)
		if err != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"Internal Server Error"})
			return
		}
		err = s.DB.ResetPassword(user.ID, string(hashedPassword))
		if err != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"Internal Server Error"})
			return
		}
		response.JSON(c, "Password Changed Successfully", http.StatusOK, nil, nil)
	}
}
