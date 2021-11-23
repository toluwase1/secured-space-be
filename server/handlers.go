package server

import (
	"github.com/decadevs/rentals-api/models"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func GetUserFromContext(c *gin.Context) (*models.User, error) {
	if userI, exists := c.Get("user"); exists {
		if user, ok := userI.(*models.User); ok {
			return user, nil
		}
		return nil, errors.New("User is not logged in")
	}
	return nil, errors.New("user is not logged in")

}
