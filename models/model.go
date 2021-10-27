package models

import "time"

type Models struct {
	ID        string    `gorm:"type:uuid;primary_key"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
