package models

// User holds a user details
type User struct {
	Models
	FirstName            string      `json:"first_name" gorm:"not null" binding:"required" form:"first_name"`
	LastName             string      `json:"last_name" gorm:"not null" binding:"required" form:"last_name"`
	Phone1               string      `json:"phone" gorm:"not null" binding:"required" form:"phone1"`
	Phone2               string      `json:"phone_2,omitempty" form:"phone2"`
	Email                string      `json:"email" gorm:"not null" binding:"required,email" form:"email"`
	Address              string      `json:"address,omitempty" gorm:"not null" form:"address"`
	HashedPassword       string      `json:"-,omitempty" gorm:"not null"`
	Password             string      `json:"password" gorm:"-" binding:"required" form:"password"`
	ConfirmPassword      string      `json:"confirm_password" gorm:"-" form:"confirm_password"`
	BookmarkedApartments []Apartment `gorm:"many2many:bookmarked_apartments" json:"bookmarked_apartments,omitempty"`
	Image                string      `json:"image,omitempty"`
	RoleID               string      `json:"role_id"`
	Role                 Role
}
type Images struct {
	Models
	ApartmentID string    `json:"apartment_id"`
	Apartment   Apartment
	URL         string    `json:"url"`
	Name        string    `json:"name"`
}
type UpdateUser struct {
	FirstName string `json:"first_name" binding:"required" form:"first_name"`
	LastName  string `json:"last_name" binding:"required" form:"last_name"`
	Phone1    string `json:"phone" binding:"required" form:"phone1"`
	Phone2    string `json:"phone_2" form:"phone2"`
	Email     string `json:"email" binding:"required,email" form:"email"`
	Address   string `json:"address"  form:"address"`
}

//https://arp-rental.s3.eu-west-3.amazonaws.com/profile_pics/61940fe40c89f633b3364156.png
