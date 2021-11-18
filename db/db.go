package db

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/decadevs/rentals-api/models"
	"mime/multipart"
)

// DB provides access to the different db
type DB interface {
	CreateUser(user *models.User) (*models.User, error)
	FindUserByUsername(username string) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
	UpdateUser(id string, update *models.UpdateUser) error
	AddToBlackList(blacklist *models.Blacklist) error
	TokenInBlacklist(token *string) bool
	FindUserByPhone(phone string) (*models.User, error)
	FindAllUsersExcept(except string) ([]models.User, error)
	CreateApartment(apartment *models.Apartment) error
	DeleteApartment(ID, userID string) error
	UpdateApartment(apartment *models.Apartment, apartmentID string) error
	GetApartmentByID(apartmentID string) (*models.Apartment, error)
	SaveBookmarkApartment(bookmarkApartment *models.BookmarkApartment) error
	CheckApartmentInBookmarkApartment(userID, apartmentID string) bool
	RemoveBookmarkedApartment(bookmarkApartment *models.BookmarkApartment) error
	GetBookmarkedApartments(userID string) ([]models.Apartment, error)
	GetUsersApartments(userId string) ([]models.Apartment, error)
	GetAllInteriorFeatures() ([]models.InteriorFeature, error)
	GetAllExteriorFeatures() ([]models.ExteriorFeature, error)
	UploadFileToS3(s *session.Session, file multipart.File, fileName string, size int64) (string, error)
	ResetPassword(userID, NewPassword string) error
	SearchApartment(categoryID, location, minPrice, maxPrice, noOfRooms string) ([]models.Apartment, error)
	GetRoleByName(name string) (models.Role, error)
}

// ValidationError defines error that occur due to validation
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Message)
}
