package models

import "time"

type Models struct {
	ID        string    `gorm:"type:uuid;primaryKey;"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
