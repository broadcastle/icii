package database

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Song is the information needed to store a song.
type Song struct {
	gorm.Model
	Album          string  `json:"album"`
	Artist         string  `json:"artist"`
	Genre          string  `json:"genre"`
	Length         float64 `json:"length"`
	Location       string  `json:"location"`
	Title          string  `json:"title"`
	Year           string  `json:"year"`
	OrganizationID uint    `json:"organization_id"`
	UserID         uint    `json:"uploader"`
	Stats          Statistics
}

// Statistics has the information about a song.
type Statistics struct {
	Played []time.Time
}

// User is the information
type User struct {
	gorm.Model
	Name          string          `json:"name"`
	Email         string          `json:"email"`
	Password      string          `json:"password"`
	Organizations []*Organization `gorm:"many2many:user_organizations;" json:"organizations"`
}

// Organization ...
type Organization struct {
	gorm.Model
	Name  string  `json:"name"`
	Slug  string  `json:"slug"`
	Users []*User `gorm:"many2many:user_organizations;" json:"users"`
	Songs []Song  `json:"songs"`
}

func initIciiTables(d *gorm.DB) {
	d.AutoMigrate(&Song{})
	d.AutoMigrate(&User{})
	d.AutoMigrate(&Organization{})
}
