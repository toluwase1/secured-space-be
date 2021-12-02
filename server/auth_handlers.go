package server

import (
	"fmt"
	"github.com/decadevs/rentals-api/db"
	"github.com/decadevs/rentals-api/models"
	"github.com/decadevs/rentals-api/server/response"
	"github.com/decadevs/rentals-api/servererrors"
	"github.com/decadevs/rentals-api/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

func (s *Server) handleSignupTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, err := s.DB.GetRoleByName("tenant")
		if err != nil {
			log.Println(err)
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"Cannot find role id"})
			return
		}
		user := &models.User{
			RoleID: role.ID,
			Role:   role,
		}

		// Generates access claims and refresh claims
		accessClaims, _ := services.GenerateClaims(user.Email)
		secret := os.Getenv("JWT_SECRET")
		accToken, err := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
		if err != nil {
			log.Printf("token generation error err: %v\n", err)
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
			return
		}

		user.Token = *accToken
		fmt.Println(*accToken)

		if errs := s.decode(c, user); errs != nil {
			response.JSON(c, "Cannot decode user signup request", http.StatusBadRequest, nil, errs)
			return
		}
		HashedPassword, err := services.GenerateHashPassword(user.Password)
		user.HashedPassword = string(HashedPassword)
		if err != nil {
			log.Printf("hash password err: %v\n", err)
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
			return
		}
		_, err = s.DB.FindUserByEmail(user.Email)
		if err == nil {
			response.JSON(c, "", http.StatusNotFound, nil, []string{"user email already exists"})
			return
		}

		_, err = s.DB.CreateUser(user)
		if err != nil {
			log.Printf("create user err: %v\n", err)
			if err, ok := err.(db.ValidationError); ok {
				response.JSON(c, "", http.StatusBadRequest, nil, []string{err.Error()})
				return
			}
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
			return
		}
		response.JSON(c, "signup successful", http.StatusCreated, nil, nil)
		_, err = s.Mail.SendVerifyAccount(user.Email, fmt.Sprintf("http://localhost:3000/verify-email/%s/%s", user.ID, *accToken))
		//_, err = s.Mail.SendVerifyAccount(user.Email,fmt.Sprintf("https://securespace-ng.herokuapp.com/api/v1/verify-email/%s/%s",user.ID,*accToken))
		if err != nil {
			log.Printf("Error: %v", err.Error())
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"Email could not be sent"})
		}
		response.JSON(c, "email sent successfully", http.StatusCreated, nil, nil)
	}
}

func (s *Server) handleSignupAgent() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, err := s.DB.GetRoleByName("agent")
		if err != nil {
			log.Println(err)
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"Cannot find role id"})
			return
		}
		user := &models.User{
			RoleID: role.ID,
			Role:   role,
		}

		// Generates access claims and refresh claims
		accessClaims, _ := services.GenerateClaims(user.Email)
		secret := os.Getenv("JWT_SECRET")
		accToken, err := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
		if err != nil {
			log.Printf("token generation error err: %v\n", err)
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
			return
		}

		user.Token = *accToken

		if errs := s.decode(c, user); errs != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, errs)
			return
		}
		HashedPassword, err := services.GenerateHashPassword(user.Password)
		user.HashedPassword = string(HashedPassword)
		if err != nil {
			log.Printf("hash password err: %v\n", err)
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
			return
		}
		_, err = s.DB.FindUserByEmail(user.Email)
		if err == nil {
			response.JSON(c, "", http.StatusNotFound, nil, []string{"user email already exists"})
			return
		}
		_, err = s.DB.CreateUser(user)
		if err != nil {
			log.Printf("create user err: %v\n", err)
			if err, ok := err.(db.ValidationError); ok {
				response.JSON(c, "", http.StatusBadRequest, nil, []string{err.Error()})
				return
			}
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
			return
		}
		response.JSON(c, "signup successful", http.StatusCreated, nil, nil)

		_, err = s.Mail.SendVerifyAccount(user.Email, fmt.Sprintf("http://localhost:3000/verify-email/%s/%s", user.ID, *accToken))
		//_, err = s.Mail.SendVerifyAccount(user.Email,fmt.Sprintf("https://securespace-ng.herokuapp.com/api/v1/verify-email/%s/%s",user.ID,*accToken))
		if err != nil {
			log.Printf("Error: %v", err.Error())
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"email could not be sent"})
		}
		response.JSON(c, "email sent successfully", http.StatusCreated, nil, nil)
	}
}

func (s *Server) handleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := &models.User{}
		loginRequest := &struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}{}

		if errs := s.decode(c, loginRequest); errs != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, errs)
			fmt.Println("here")
			return
		}
		// Check if the user with that email exists
		user, err := s.DB.FindUserByEmail(loginRequest.Email)

		if err != nil {
			if inactiveErr, ok := err.(servererrors.InActiveUserError); ok {
				response.JSON(c, "", http.StatusBadRequest, nil, []string{inactiveErr.Error()})
				return
			}
			log.Printf("No user: %v\n", err)
			response.JSON(c, "", http.StatusUnauthorized, nil, []string{"user not found"})
			return
		}

		if user.IsActive == false {
			log.Printf("No user: %v\n", err)
			response.JSON(c, "", http.StatusUnauthorized, nil, []string{"email not verified"})
			return
		}

		err = services.CompareHashAndPassword([]byte(loginRequest.Password), user.HashedPassword)
		if err != nil {
			log.Printf("passwords do not match %v\n", err)
			response.JSON(c, "", http.StatusUnauthorized, nil, []string{"email or password incorrect"})
			return
		}

		// Generates access claims and refresh claims
		accessClaims, refreshClaims := services.GenerateClaims(user.Email)

		secret := os.Getenv("JWT_SECRET")
		accToken, err := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
		if err != nil {
			log.Printf("token generation error err: %v\n", err)
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
			return
		}

		refreshToken, err := services.GenerateToken(jwt.SigningMethodHS256, refreshClaims, &secret)
		if err != nil {
			log.Printf("token generation error err: %v\n", err)
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
			return
		}

		response.JSON(c, "login successful", http.StatusOK, gin.H{
			"user":          user,
			"access_token":  *accToken,
			"refresh_token": *refreshToken,
		}, nil)
	}
}

func (s *Server) handleLogout() gin.HandlerFunc {
	return func(c *gin.Context) {

		if tokenI, exists := c.Get("access_token"); exists {
			if userI, exists := c.Get("user"); exists {
				if user, ok := userI.(*models.User); ok {
					if accessToken, ok := tokenI.(string); ok {

						rt := &struct {
							RefreshToken string `json:"refresh_token,omitempty" binding:"required"`
						}{}

						if err := c.ShouldBindJSON(rt); err != nil {
							log.Printf("no refresh token in request body: %v\n", err)
							response.JSON(c, "", http.StatusBadRequest, nil, []string{"unauthorized"})
							return
						}

						accBlacklist := &models.Blacklist{
							Email:     user.Email,
							CreatedAt: time.Now(),
							Token:     accessToken,
						}

						err := s.DB.AddToBlackList(accBlacklist)
						if err != nil {
							log.Printf("can't add access token to blacklist: %v\n", err)
							response.JSON(c, "logout failed", http.StatusInternalServerError, nil, []string{"couldn't revoke access token"})
							return
						}

						refreshBlacklist := &models.Blacklist{
							Email:     user.Email,
							CreatedAt: time.Now(),
							Token:     rt.RefreshToken,
						}

						err = s.DB.AddToBlackList(refreshBlacklist)
						if err != nil {
							log.Printf("can't add refresh token to blacklist: %v\n", err)
							response.JSON(c, "logout failed", http.StatusInternalServerError, nil, []string{"couldn't revoke refresh token"})
							return
						}
						response.JSON(c, "logout successful", http.StatusOK, nil, nil)
						return
					}
				}
			}
		}
		log.Printf("can't get info from context\n")
		response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
		return
	}
}
