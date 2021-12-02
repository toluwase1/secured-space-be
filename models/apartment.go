package models

type ApartmentStatus bool

const (
	Available    ApartmentStatus = true
	NotAvailable ApartmentStatus = false
)

type Apartment struct {
	Models
	UserID          string `json:"user_id" gorm:"foreignkey:User(id)"`
	User            User
	Title           string            `json:"title" gorm:"not null"`
	CategoryID      string            `json:"category_id"  gorm:"foreignkey: categories(id)"`
	Description     string            `json:"description"  gorm:"not null"`
	Price           int               `json:"price"  gorm:"not null"`
	NoOfRooms       int               `json:"no_of_rooms"  gorm:"not null"`
	Furnished       bool              `json:"furnished"  gorm:"not null"`
	Location        string            `json:"location" gorm:"not null"`
	ApartmentStatus ApartmentStatus   `json:"apartment_status" gorm:"not null default:true"`
	Images          []Images          `json:"images" gorm:"not null"`
	Interiors       []InteriorFeature `json:"interior" gorm:"not null; many2many:apartment_interiors"`
	Exteriors       []ExteriorFeature `json:"exterior" gorm:"not null; many2many:apartment_exteriors"`
}
