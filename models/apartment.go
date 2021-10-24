package models

type Apartment struct {
	Id 				string     `gorm:"type:uuid;primaryKey; default:uuid_generate_v4();size:100; not null"`
	UserId      	string     `json:"user_id" gorm:"foreignkey:User(id)"`
	Title         	string	   `json:"title" gorm:"not null";<-`
	CategoryId   	string	   `json:"category_id"  gorm:"foreignkey: categories(id)"`
	Description   	string     `json:"description"  gorm:"not null"`
	Price         	int        `json:"price"  gorm:"not null"`
	NoOfRooms    	int        `json:"no_of_rooms"  gorm:"not null"`
	Furnished     	bool       `json:"furnished"  gorm:"not null"`
	Location      	string     `json:"location" gorm:"not null"`
	Status        	bool       `json:"status" gorm:"not null"`
	Images        	[]string   `json:"images" gorm:"not null"`

}