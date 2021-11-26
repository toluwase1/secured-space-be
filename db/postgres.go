package db

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/decadevs/rentals-api/models"
	"github.com/google/uuid"
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
	var dsn string
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		dsn = fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v", DBHost, DBUser, DBPass, DBName, DBPort, DBMode, DBTimeZone)
	} else {
		dsn = databaseUrl
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	postgresDB.DB = db
}
func (postgresDB *PostgresDB) PopulateTables() {
	err := postgresDB.DB.AutoMigrate(&models.User{}, &models.Role{}, &models.Apartment{}, &models.Images{}, &models.InteriorFeature{}, &models.ExteriorFeature{}, &models.Category{}, &models.Blacklist{})
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

	categories := []models.Category{{Name: "bungalow"}, {Name: "townhouse"}, {Name: "terraced-houses"}, {Name: "penthouse"}, {Name: "semi-detached"}, {Name: "maisonette"}, {Name: "duplex"}}
	postgresDB.DB.Create(&categories)

	interiorFeatures := []models.InteriorFeature{
		{ID: uuid.NewString(), Name: "adsl"},
		{ID: uuid.NewString(), Name: "barbecue"},
		{ID: uuid.NewString(), Name: "panel door"},
		{ID: uuid.NewString(), Name: "ceramic floor"},
		{ID: uuid.NewString(), Name: "balcony"},
		{ID: uuid.NewString(), Name: "alarm"},
		{ID: uuid.NewString(), Name: "laminate"},
		{ID: uuid.NewString(), Name: "blinds"},
		{ID: uuid.NewString(), Name: "sauna"},
		{ID: uuid.NewString(), Name: "laundry room"},
		{ID: uuid.NewString(), Name: "video intercom"},
		{ID: uuid.NewString(), Name: "shower"},
		{ID: uuid.NewString(), Name: "dressing room"},
		{ID: uuid.NewString(), Name: "satin plaster"},
		{ID: uuid.NewString(), Name: "wallpaper"},
	}
	postgresDB.DB.Create(&interiorFeatures)

	exteriorFeatures := []models.ExteriorFeature{
		{ID: uuid.NewString(), Name: "car park"},
		{ID: uuid.NewString(), Name: "elevator"},
		{ID: uuid.NewString(), Name: "tennis court"},
		{ID: uuid.NewString(), Name: "gym"},
		{ID: uuid.NewString(), Name: "garden"},
		{ID: uuid.NewString(), Name: "basketball court"},
		{ID: uuid.NewString(), Name: "thermal insulation"},
		{ID: uuid.NewString(), Name: "market"},
		{ID: uuid.NewString(), Name: "security"},
		{ID: uuid.NewString(), Name: "pvc"},
		{ID: uuid.NewString(), Name: "generator"},
	}
	postgresDB.DB.Create(&exteriorFeatures)
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

func (postgresDB *PostgresDB) FindUserByID(userID string) (*models.User, error) {
	var user *models.User
	err := postgresDB.DB.Where("id = ?", userID).First(&user).Error
	return user, err
}


func (postgresDB *PostgresDB)SetUserToActive(userID string)  error{
	var user *models.User
	err := postgresDB.DB.Model(&user).Where("id = ?", userID).Update("is_active", true).Error
	return err
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
	result := postgresDB.DB.Where("token = ?", token).Find(&models.Blacklist{})
	return result.Error != nil
}
func (postgresDB *PostgresDB) FindUserByPhone(phone string) (*models.User, error) {
	return nil, nil
}
func (postgresDB *PostgresDB) FindAllUsersExcept(except string) ([]models.User, error) {
	return nil, nil
}

func (postgresDB *PostgresDB) GetUsersApartments(userId string) ([]models.Apartment, error) {
	var Apartments []models.Apartment

	result := postgresDB.DB.Preload("Images").Preload("User").Where("user_id=?", userId).Find(&Apartments)

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
	db := postgresDB.DB.Table("bookmarked_apartments").Create(&bookmarkApartment)
	return db.Error
}

func (postgresDB *PostgresDB) CheckApartmentInBookmarkApartment(userID, apartmentID string) bool {
	result := postgresDB.DB.Table("bookmarked_apartments").Where("user_id = ? AND apartment_id = ?", userID, apartmentID).First(&models.BookmarkApartment{})
	return result.RowsAffected == 1
}
func (postgresDB *PostgresDB) UpdateApartment(apartment *models.Apartment, apartmentID string) error {
	result := postgresDB.DB.Model(models.Apartment{}).Where("id = ?", apartmentID).Updates(apartment)
	return result.Error
}

func (postgresDB *PostgresDB) RemoveBookmarkedApartment(bookmarkApartment *models.BookmarkApartment) error {
	result := postgresDB.DB.Table("bookmarked_apartments").
		Where("user_id = ? AND apartment_id = ?", bookmarkApartment.UserID, bookmarkApartment.ApartmentID).
		Delete(&models.BookmarkApartment{})
	return result.Error
}

func (postgresDB *PostgresDB) GetBookmarkedApartments(userID string) ([]models.Apartment, error) {
	user := &models.User{}
	result := postgresDB.DB.Preload("BookmarkedApartments.Images").Where("id = ?", userID).Find(&user)
	return user.BookmarkedApartments, result.Error
}

func (postgresDB *PostgresDB) GetAllCategory() ([]models.Category, error) {
	categories := []models.Category{}
	if err := postgresDB.DB.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (postgresDB *PostgresDB) GetAllInteriorFeatures() ([]models.InteriorFeature, error) {
	interiorFeatures := []models.InteriorFeature{}
	if err := postgresDB.DB.Find(&interiorFeatures).Error; err != nil {
		return nil, err
	}
	return interiorFeatures, nil
}

func (postgresDB *PostgresDB) GetAllExteriorFeatures() ([]models.ExteriorFeature, error) {
	exteriorFeatures := []models.ExteriorFeature{}
	if err := postgresDB.DB.Find(&exteriorFeatures).Error; err != nil {
		return nil, err
	}
	return exteriorFeatures, nil
}

func (p *PostgresDB) UploadFileToS3(s *session.Session, file multipart.File, fileName string, size int64) (string, error) {
	// get the file size and read
	// the file content into a buffer
	buffer := make([]byte, size)
	file.Read(buffer)
	// config settings: this is where you choose the bucket,
	// filename, content-type and storage class of the file
	// you're uploading
	url := "https://s3-eu-west-3.amazonaws.com/arp-rental/" + fileName
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
	return url, err
}

func (postgresDB *PostgresDB) ResetPassword(userID, NewPassword string) error {
	result := postgresDB.DB.Model(models.User{}).Where("id = ?", userID).Update("hashed_password", NewPassword)
	return result.Error
}

func (postgresDB *PostgresDB) SearchApartment(categoryID, location, minPrice, maxPrice, noOfRooms string) ([]models.Apartment, error) {
	var apartments []models.Apartment
	stm := ""
	if minPrice == "" && maxPrice == "" {
		stm = fmt.Sprintf("(category_id LIKE '%%%s%%' AND location LIKE '%%%s%%')", categoryID, location)
	} else if noOfRooms != "" {
		stm = fmt.Sprintf("(no_of_rooms <= %s OR (price >= %s AND price <= %s)) AND (category_id LIKE '%%%s%%' AND location LIKE '%%%s%%')", noOfRooms, minPrice, maxPrice, categoryID, location)
	} else if minPrice != "" {
		stm = fmt.Sprintf("(( price >= %s)) AND (category_id LIKE '%%%s%%' AND location LIKE '%%%s%%')", minPrice, categoryID, location)
	} else if maxPrice != "" {
		stm = fmt.Sprintf("(( price <= %s)) AND (category_id LIKE '%%%s%%' AND location LIKE '%%%s%%')", maxPrice, categoryID, location)
	} else if location != "0" {
		stm = fmt.Sprintf("(location LIKE '%%%s%%')", location)
	} else if categoryID != "0" {
		stm = fmt.Sprintf("(category_id LIKE '%%%s%%')", categoryID)
	} else if categoryID == "0" && location == "0" && minPrice == "" && maxPrice == "" {
		result := postgresDB.DB.Preload("Images").Find(&apartments)
		return apartments, result.Error
	} else {
		stm = fmt.Sprintf("((price >= %s AND price <= %s)) AND (category_id LIKE '%%%s%%' AND location LIKE '%%%s%%')", minPrice, maxPrice, categoryID, location)
	}
	result := postgresDB.DB.Preload("Images").Where(stm).Find(&apartments)
	return apartments, result.Error
}

func (postgersDB *PostgresDB) ApartmentDetails(apartmentID string) (*models.Apartment, error) {
	var apart *models.Apartment
	result := postgersDB.DB.Preload("Images").Preload("User").Preload("Exteriors").Preload("Interiors").Where("id = ?", apartmentID).Find(&apart)
	return apart, result.Error
}

func (postgresDB *PostgresDB) GetRoleByName(name string) (models.Role, error) {
	var role models.Role
	err := postgresDB.DB.Where("title = ?", name).First(&role).Error
	return role, err
}

func (PostgresDB *PostgresDB) GetApartmentByCategory(categoryID string) []models.Apartment {
	var apartments []models.Apartment
	PostgresDB.DB.Preload("Images").Where("category_id = ?", categoryID).Find(&apartments)
	return apartments
}

func (PostgresDB *PostgresDB) GetAllCategories() []models.Category {
	var categories []models.Category
	PostgresDB.DB.Find(&categories)
	return categories
}