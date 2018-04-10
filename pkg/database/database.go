package database

import "github.com/jinzhu/gorm"

// Song is the information needed to store a song.
type Song struct {
	gorm.Model
	Name     string  `json:"name"`
	Artist   string  `json:"artist"`
	Album    string  `json:"album"`
	Length   float64 `json:"length"`
	Location string  `json:"location"`
}

// User is the information
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	API      string `json:"api_key"`
}

// Organization ...
type Organization struct {
	gorm.Model
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func initIciiTables(d *gorm.DB) {
	d.AutoMigrate(&Song{})
	d.AutoMigrate(&User{})
}
