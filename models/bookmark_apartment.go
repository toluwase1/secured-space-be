package models

type BookmarkApartment struct {
	Models
	UserID      string `json:"user_id"`
	ApartmentID string `json:"apartment_id"`
}
