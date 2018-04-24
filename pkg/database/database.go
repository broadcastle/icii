package database

import (
	"github.com/jinzhu/gorm"
)

// Track is the information needed to store a audio file.
type Track struct {
	gorm.Model
	Album     string  `json:"album"`
	Artist    string  `json:"artist"`
	Location  string  `json:"location"`
	Title     string  `json:"title"`
	StationID uint    `json:"organization_id"`
	UserID    uint    `json:"uploader"`
	Year      string  `json:"year"`
	Length    float64 `json:"length"`
	Genre     string  `json:"genre"`
	Tags      []Tag   `json:"tags"`
}

// Tag is for the tags that would be assigned to a audio file.
type Tag struct {
	gorm.Model
	Text    string `json:"text"`
	TrackID uint   `json:"audio_id"`
}

// User is the information about the user for icii.
type User struct {
	gorm.Model
	Name     string     `json:"name"`
	Email    string     `json:"email"`
	Password string     `json:"password"`
	Stations []*Station `gorm:"many2many:user_stations;" json:"organizations"`
}

// A Station holds the audio tracks and users.
type Station struct {
	gorm.Model
	Name  string  `json:"name"`
	Slug  string  `json:"slug"`
	Users []*User `gorm:"many2many:user_stations;" json:"users"`
	Track []Track `json:"audio"`
}

// UserPermission keeps track of what permission are allowed for a user.
type UserPermission struct {
	gorm.Model
	UserID    uint `json:"user_id"`
	StationID uint `json:"station_id"`

	TrackAdd    bool `json:"track_add"`
	TrackEdit   bool `json:"track_edit"`
	TrackRemove bool `json:"track_remove"`

	UserAdd    bool `json:"user_add"`
	UserEdit   bool `json:"user_edit"`
	UserRemove bool `json:"user_remove"`

	StreamAdd    bool `json:"stream_add"`
	StreamEdit   bool `json:"stream_edit"`
	StreamRemove bool `json:"stream_remove"`

	StationAdd    bool `json:"station_add"`
	StationEdit   bool `json:"station_edit"`
	StationRemove bool `json:"station_remove"`
}

// Initialize the database tables.
func initIciiTables(d *gorm.DB) {
	d.AutoMigrate(&Station{})
	d.AutoMigrate(&Tag{})
	d.AutoMigrate(&Track{})
	d.AutoMigrate(&User{})
	d.AutoMigrate(&UserPermission{})
}
