module github.com/decadevs/rentals-api

// +heroku goVersion go1.17
go 1.17

require (
	github.com/aws/aws-sdk-go v1.42.9
	github.com/brianvoe/gofakeit/v6 v6.9.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.4
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/go-playground/validator/v10 v10.9.0
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.3.0
	github.com/joho/godotenv v1.4.0
	github.com/mailgun/mailgun-go/v4 v4.6.0
	github.com/pkg/errors v0.9.1
	github.com/pusher/pusher-http-go v4.0.1+incompatible
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20211117183948-ae814b36b871
	gorm.io/driver/postgres v1.2.2
	gorm.io/gorm v1.22.3
)
