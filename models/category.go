package models

import "time"

type Category struct {
	ID        string `json:"id"`
	Name      string `json:"Name"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
