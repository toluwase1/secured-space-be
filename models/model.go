package models

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type Models struct {
	ID        string    `gorm:"type:uuid;primaryKey;"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (u *Models) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	if u.ID == "" {
		err = errors.New("can't save invalid data")
	}
	return
}
