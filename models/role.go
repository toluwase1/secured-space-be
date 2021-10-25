package models

import "gorm.io/gorm"

type Role struct {
	Title string
	gorm.Model
}
