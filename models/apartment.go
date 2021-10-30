package models

type ApartmentStatus bool

const (
	Available    ApartmentStatus = true
	NotAvailable ApartmentStatus = false
)

type Apartment struct {
	ID              string          `sql:"type:uuid; default:uuid_generate_v4();size:100; not null"`
	UserID          string          `json:"user_id" gorm:"foreignkey:User(id)"`
	Title           string          `json:"title" gorm:"not null"`
	CategoryID      string          `json:"category_id"  gorm:"foreignkey: categories(id)"`
	Description     string          `json:"description"  gorm:"not null"`
	Price           int             `json:"price"  gorm:"not null"`
	NoOfRooms       int             `json:"no_of_rooms"  gorm:"not null"`
	Furnished       bool            `json:"furnished"  gorm:"not null"`
	Location        string          `json:"location" gorm:"not null"`
	ApartmentStatus ApartmentStatus `json:"apartment_status" gorm:"not null"`
	//Images        	[]Image    `json:"images" gorm:"not null"`
	//Interiors		[]Interior  `json:"interior" gorm:"not null; many2many:apartment_interior"`
	//Exteriors     []Exterior   `json: "exterior" gorm: "not null; many2many:apartment_exterior"`
}
