package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Title string `json:"title"`
}
