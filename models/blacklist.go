package models

import "time"

//Blacklist helps us blacklist tokens
type Blacklist struct {
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}
