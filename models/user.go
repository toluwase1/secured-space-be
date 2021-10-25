package models

import "time"

// User holds a user details
type User struct {
	Models
	FirstName       string `json:"first_name" gorm:"type:not null" binding:"required" form:"first_name"`
	LastName        string `json:"last_name" gorm:"type:not null" binding:"required" form:"last_name"`
	Phone1          string `json:"phone" gorm:"type:not null" binding:"required" form:"phone1"`
	Phone2          string `json:"phone_2" form:"phone2"`
	Email           string `json:"email" gorm:"type:not null" binding:"required,email" form:"email"`
	Address         string `json:"address" gorm:"type:not null" binding:"required" form:"address"`
	HashedPassword  []byte `json:"-" gorm:"type:not null"`
	Password        string `json:"password" gorm:"-" binding:"required" form:"password"`
	ConfirmPassword string `json:"confirm_password" gorm:"-" form:"confirm_password"`
	Image           string `json:"image,omitempty"`
	RoleID          int    `json:"role_id"`
	Role            Role
}
type Images struct {
	ID          string `json:"id"`
	ApartmentID string `json:"apartment_id"`
	URL         string
	Name        string `json:"name"`
	CreateAt    time.Time
	UpdateAt    time.Time
}
