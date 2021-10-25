package models

type Feature struct {
	ID        string    `json:"id"`
	//Apartment Apartment `json:"apartment_id" gorm:"foreignKey:ID"`
	Interior  string    `json:"interior"`
	Exterior  string    `json:"exterior"`
}
