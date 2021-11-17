package services

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const AccessTokenValidity = time.Hour * 24 * 30
const RefreshTokenValidity = time.Hour * 24

// GetTokenFromHeader returns the token string in the authorization header
func GetTokenFromHeader(c *gin.Context) string {
	authHeader := c.Request.Header.Get("Authorization")
	if len(authHeader) > 8 {
		return authHeader[7:]
	}
	return ""
}

// verifyAccessToken verifies a token
func verifyToken(tokenString *string, claims jwt.MapClaims, secret *string) (*jwt.Token, error) {
	parser := &jwt.Parser{SkipClaimsValidation: true}
	return parser.ParseWithClaims(*tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(*secret), nil
	})
}

// AuthorizeToken check if a refresh token is valid
func AuthorizeToken(token *string, secret *string) (*jwt.Token, jwt.MapClaims, error) {
	if token != nil && *token != "" && secret != nil && *secret != "" {
		claims := jwt.MapClaims{}
		token, err := verifyToken(token, claims, secret)
		if err != nil {
			return nil, nil, err
		}
		return token, claims, nil
	}
	return nil, nil, fmt.Errorf("empty token or secret")
}

// GenerateToken generates only an access token
func GenerateToken(signMethod *jwt.SigningMethodHMAC, claims jwt.MapClaims, secret *string) (*string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(signMethod, claims)
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(*secret))
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func CompareHashAndPassword(password []byte, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), password)
}

func GenerateHashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func GenerateClaims(email string) (jwt.MapClaims, jwt.MapClaims) {
	accessClaims := jwt.MapClaims{
		"user_email": email,
		"exp":        time.Now().Add(AccessTokenValidity).Unix(),
	}
	refreshClaims := jwt.MapClaims{
		"exp": time.Now().Add(RefreshTokenValidity).Unix(),
		"sub": 1,
	}

	return accessClaims, refreshClaims
}
