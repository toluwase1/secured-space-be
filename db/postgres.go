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

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v", DBHost, DBUser, DBPass,DBName, DBPort, DBMode, DBTimeZone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	postgresDB.DB = db
	//err = postgresDB.DB.AutoMigrate(&models.Role{})
	//roles := []models.Role{{Title: "tenant"},{Title: "agent"}}
	//postgresDB.DB.Create(roles)
	//if err != nil {
	//	return
	//}
	err = postgresDB.DB.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Printf("err %v:",err.Error())

	}

}

var user *models.User

func (postgresDB *PostgresDB) CreateUser(user *models.User) (*models.User, error) {
	result := postgresDB.DB.Create(&user)
	return user, result.Error
}
func (postgresDB *PostgresDB) FindUserByUsername(username string) (*models.User, error) {
	result := postgresDB.DB.Where("username = ?",username).First(&user)
	return user, result.Error

}
func (postgresDB *PostgresDB) FindUserByEmail(email string) (*models.User, error) {
	result := postgresDB.DB.Where("email = ?",email).First(&user)
	return user, result.Error
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
