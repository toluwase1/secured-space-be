package db

import (
	"fmt"
	"github.com/decadevs/rentals-api/models"
)

// DB provides access to the different db
type DB interface {
	CreateUser(user *models.User) (*models.User, error)
	FindUserByUsername(username string) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	AddToBlackList(blacklist *models.Blacklist) error
	TokenInBlacklist(token *string) bool
	FindUserByPhone(phone string) (*models.User, error)
	UpdateApartment(apartment *models.Apartment, apartmentID string) error
	SaveBookmarkApartment(bookmarkApartment *models.BookmarkApartment) error
	CheckApartmentInBookmarkApartment(userID, apartmentID string) bool

}

// ValidationError defines error that occur due to validation
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Message)
}
