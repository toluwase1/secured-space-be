package db

import (
	"fmt"
	"github.com/decadevs/rentals-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
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
		log.Panicln(err.Error())

	}

}

func (postgresDB *PostgresDB) CreateUser(user *models.User) (*models.User, error) {
	return nil, nil
}
func (postgresDB *PostgresDB) FindUserByUsername(username string) (*models.User, error) {
	return nil, nil
}
func (postgresDB *PostgresDB) FindUserByEmail(email string) (*models.User, error) {
	var user *models.User
	userEmail := postgresDB.DB.Where("email = ?", email).Preload("Role").First(&user)
	return user, userEmail.Error
}
func (postgresDB *PostgresDB) UpdateUser(user *models.User) error {
	return nil
}
func (postgresDB *PostgresDB) AddToBlackList(blacklist *models.Blacklist) error {
	return nil
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
