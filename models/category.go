package models

import "time"

type Category struct {
	ID        string `json:"id"`
	Name      string `json:"shop"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
