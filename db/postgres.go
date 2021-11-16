package db

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/decadevs/rentals-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

// PostgresDB implements the DB interface
type PostgresDB struct {
	DB *gorm.DB
}

// Init sets up the mongodb instance
func (postgresDB *PostgresDB) Init() {
	// Database Variables
	DBUser := os.Getenv("DB_USER")
	DBPass := os.Getenv("DB_PASS")
	DBHost := os.Getenv("DB_HOST")
	DBName := os.Getenv("DB_NAME")
	DBPort := os.Getenv("DB_PORT")
	DBTimeZone := os.Getenv("DB_TIMEZONE")
	DBMode := os.Getenv("DB_MODE")

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v", DBHost, DBUser, DBPass, DBName, DBPort, DBMode, DBTimeZone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	postgresDB.DB = db

	err = postgresDB.DB.AutoMigrate(&models.User{}, &models.Role{}, &models.Apartment{}, &models.Images{}, &models.InteriorFeature{}, &models.ExteriorFeature{}, &models.Category{})
	if err != nil {
		log.Println("unable to migrate database.", err.Error())
	}

	err = postgresDB.DB.Create(&models.Role{Title: "tenant"}).Error
	if err != nil {
		log.Println("unable to create role.", err.Error())
	}
	err = postgresDB.DB.Create(&models.Role{Title: "agent"}).Error
	if err != nil {
		log.Println("unable to create role.", err.Error())
	}

}

func (postgresDB *PostgresDB) CreateUser(user *models.User) (*models.User, error) {
	err := postgresDB.DB.Create(user).Error
	return nil, err
}
func (postgresDB *PostgresDB) FindUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := postgresDB.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}
func (postgresDB *PostgresDB) FindUserByEmail(email string) (*models.User, error) {
	var user *models.User
	userEmail := postgresDB.DB.Where("email = ?", email).Preload("Role").First(&user)
	return user, userEmail.Error
}
func (postgresDB *PostgresDB) UpdateUser(id string, update *models.UpdateUser) error {
	result :=
		postgresDB.DB.Model(models.User{}).
			Where("id = ?", id).
			Updates(
				models.User{
					FirstName: update.FirstName,
					LastName:  update.LastName,
					Phone1:    update.Phone1,
					Phone2:    update.Phone2,
					Address:   update.Address,
					Email:     update.Email,
				},
			)
	return result.Error
}
func (postgresDB *PostgresDB) AddToBlackList(blacklist *models.Blacklist) error {
	result := postgresDB.DB.Create(blacklist)
	return result.Error
}
func (postgresDB *PostgresDB) TokenInBlacklist(token *string) bool {
	return false
}
func (postgresDB *PostgresDB) FindUserByPhone(phone string) (*models.User, error) {
	return nil, nil
}
func (postgresDB *PostgresDB) FindAllUsersExcept(except string) ([]models.User, error) {
	return nil, nil
}

func (postgresDB *PostgresDB) GetUsersApartments(userId string) ([]models.Apartment, error) {
	var Apartments []models.Apartment

	result := postgresDB.DB.Where("user_id=?", userId).Find(&Apartments)

	return Apartments, result.Error
}

func (postgresDB *PostgresDB) CreateApartment(apartment *models.Apartment) error {
	err := postgresDB.DB.Create(&apartment).Error
	return err
}

func (postgresDB *PostgresDB) DeleteApartment(ID, userID string) error {
	result := postgresDB.DB.Where("id = ? AND user_id = ?", ID, userID).Delete(&models.Apartment{})
	return result.Error
}
func (postgresDB *PostgresDB) SaveBookmarkApartment(bookmarkApartment *models.BookmarkApartment) error {
	db := postgresDB.DB.Create(&bookmarkApartment)
	return db.Error
}

func (postgresDB *PostgresDB) CheckApartmentInBookmarkApartment(userID, apartmentID string) bool {
	result := postgresDB.DB.Where("user_id = ? AND apartment_id = ?", userID, apartmentID).First(&models.BookmarkApartment{})
	return result.RowsAffected == 1
}
func (postgresDB *PostgresDB) UpdateApartment(apartment *models.Apartment, apartmentID string) error {
	result := postgresDB.DB.Model(models.Apartment{}).Where("id = ?", apartmentID).Updates(apartment)
	return result.Error
}

func (postgresDB *PostgresDB) RemoveBookmarkedApartment(bookmarkApartment *models.BookmarkApartment) error {
	result := postgresDB.DB.
		Where("user_id = ? AND apartment_id = ?", bookmarkApartment.UserID, bookmarkApartment.ApartmentID).
		Delete(&models.BookmarkApartment{})
	return result.Error
}

func (postgresDB *PostgresDB) GetBookmarkedApartments(userID string) ([]models.Apartment, error) {
	user := &models.User{}
	result := postgresDB.DB.Preload("BookmarkApartment").Where("id = ?", userID).Find(&user)
	return user.BookmarkedApartments, result.Error
}

func (p *PostgresDB) UploadFileToS3(s *session.Session, file multipart.File, fileName string, size int64) error {
	// get the file size and read
	// the file content into a buffer
	buffer := make([]byte, size)
	file.Read(buffer)
	// config settings: this is where you choose the bucket,
	// filename, content-type and storage class of the file
	// you're uploading

	_, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:                  aws.String(fileName),
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(int64(size)),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("INTELLIGENT_TIERING"),
	})
	return err
}
func (postgresDB *PostgresDB) ResetPassword(userID, NewPassword string) error {
	result := postgresDB.DB.Model(models.User{}).Where("id = ?", userID).Update("hashed_password", NewPassword)
	return result.Error
}

func (postgresDB *PostgresDB) SearchApartment(categoryID, location, minPrice, maxPrice, noOfRooms string) ([]models.Apartment, error) {
	var apartments []models.Apartment
	stm := ""
	if minPrice == "" {
		stm = fmt.Sprintf("((price = %s)) AND (category_id LIKE '%%%s%%' AND location LIKE '%%%s%%')", maxPrice, categoryID, location)
	} else if noOfRooms != "" {
		stm = fmt.Sprintf("(no_of_rooms <= %s OR (price >= %s AND price <= %s)) AND (category_id LIKE '%%%s%%' AND location LIKE '%%%s%%')", noOfRooms, minPrice, maxPrice, categoryID, location)
	} else {
		stm = fmt.Sprintf("((price >= %s AND price <= %s)) AND (category_id LIKE '%%%s%%' AND location LIKE '%%%s%%')", minPrice, maxPrice, categoryID, location)

	}
	result := postgresDB.DB.Preload("Images").Where(stm).Find(&apartments)
	return apartments, result.Error
}

func (postgresDB *PostgresDB) GetRoleByName(name string) (models.Role, error) {
	var role models.Role
	err := postgresDB.DB.Where("title = ?", name).First(&role).Error
	return role, err
}