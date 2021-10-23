package db

import (
	"log"

	"github.com/decadevs/rentals-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresDB implements the DB interface
type PostgresDB struct {
	DB *gorm.DB
}

// Init sets up the mongodb instance
func (postgresDB *PostgresDB) Init() {
	dsn := "host=192.168.0.1 user=postgres password=postgres dbname=rental port=5432 sslmode=disable TimeZone=Africa/Lagos"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	postgresDB.DB = db
}

func (postgresDB *PostgresDB) CreateUser(user *models.User) (*models.User, error) {
	return nil, nil
}
func (postgresDB *PostgresDB) FindUserByUsername(username string) (*models.User, error) {
	return nil, nil
}
func (postgresDB *PostgresDB) FindUserByEmail(email string) (*models.User, error) {
	return nil, nil
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
