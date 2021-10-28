package server

import (
	"github.com/decadevs/rentals-api/servererrors"
	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
)

// decode decodes the body of c into v
func (s *Server) decode(c *gin.Context, v interface{}) []string {
	if err := c.ShouldBindJSON(v); err != nil {
		errs := []string{}
		verr, ok := err.(validator.ValidationErrors)
		if ok {
			for _, fieldErr := range verr {
				errs = append(errs, servererrors.NewFieldError(fieldErr).String())
			}
		} else {
			errs = append(errs, "internal server error")
		}
		return errs
	}
	return nil
}
