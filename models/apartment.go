package models

type ApartmentStatus bool

const (
	Available    ApartmentStatus = true
	NotAvailable ApartmentStatus = false
)

type Apartment struct {
	Models
	UserID          string            `json:"user_id" gorm:"foreignkey:User(id)"`
	Title           string            `json:"title" gorm:"not null"`
	CategoryID      string            `json:"category_id"  gorm:"foreignkey: categories(id)"`
	Description     string            `json:"description"  gorm:"not null"`
	Price           int               `json:"price"  gorm:"not null"`
	NoOfRooms       int               `json:"no_of_rooms"  gorm:"not null"`
	Furnished       bool              `json:"furnished"  gorm:"not null"`
	Location        string            `json:"location" gorm:"not null"`
	ApartmentStatus ApartmentStatus   `json:"apartment_status" gorm:"not null"`
	Images          []Images          `json:"images" gorm:"not null"`
	Interior        []InteriorFeature `json:"interior" gorm:"not null; many2many:apartment_interior"`
	Exterior        []ExteriorFeature `json:"exterior" gorm:"not null; many2many:apartment_exterior"`
}
