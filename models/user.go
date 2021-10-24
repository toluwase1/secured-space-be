package models

import (
	"time"
)

// User holds a user details
type User struct {
	ID              string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	FirstName       string `json:"first_name" gorm:"type:not null" binding:"required" form:"first_name"`
	LastName        string `json:"last_name" gorm:"type:not null" binding:"required" form:"last_name"`
	Phone1          string `json:"phone" gorm:"type:not null" binding:"required" form:"phone1"`
	Phone2          string `json:"phone_2" form:"phone2"`
	Email           string `json:"email" gorm:"type:not null" binding:"required,email" form:"email"`
	Address         string `json:"address" gorm:"type:not null" binding:"required" form:"address"`
	Password        string `json:"-" gorm:"type:not null" binding:"required" form:"password"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
	Image           string `json:"image,omitempty"`
	//RoleID          Role      `json:"role_id" gorm:"foreignKey:RoleID"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
