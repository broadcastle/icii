package database

import (
	"github.com/jinzhu/gorm"
)

// Track is the information needed to store a audio file.
type Track struct {
	gorm.Model

	// Necesary information.
	Album  string `json:"album"`
	Artist string `json:"artist"`
	Title  string `json:"title"`
	Year   string `json:"year"`

	// Optional information.
	Genre string `json:"genre"`
	Tags  []Tag  `json:"tags"`

	// Calculated information.
	Length    float64 `json:"length"`
	Location  string  `json:"location"`
	StationID uint    `json:"organization_id"`
	Tempo     int     `json:"tempo"`
	UserID    uint    `json:"uploader"`
}

// Tag holds information that can be assigned to multiple tracks within a station.
type Tag struct {
	gorm.Model

	// Information that is visible to the station users.
	Text      string `json:"text"`
	TrackID   uint   `json:"audio_id"`
	StationID uint   `json:"station_id"`
}

// User is the information about the user for icii.
type User struct {
	gorm.Model

	// Only the user can change this information.
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`

	// Relationship between the users and the stations.
	Stations []*Station `gorm:"many2many:user_stations;" json:"organizations"`
}

// A Station holds the audio tracks and users.
type Station struct {
	gorm.Model

	// User filled in information.
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Public bool   `json:"public"`

	// Relationships between users and tracks.
	Users []*User `gorm:"many2many:user_stations;" json:"users"`
	Track []Track `json:"audio"`
}

// UserPermission keeps track of what permission are allowed for a user.
type UserPermission struct {
	gorm.Model

	UserID    uint `json:"user_id"`
	StationID uint `json:"station_id"`

	// Station track permisions.
	TrackAdd    bool `json:"track_add"`    // Add tracks to the station.
	TrackEdit   bool `json:"track_edit"`   // Edit a stations tracks.
	TrackRemove bool `json:"track_remove"` // Remove a stations tracks.

	// Station user permissions.
	UserAdd    bool `json:"user_add"`    // Add users to the station.
	UserEdit   bool `json:"user_edit"`   // Edit users permissions in the station.
	UserRemove bool `json:"user_remove"` // Remove user from a station.

	// Station stream permissions.
	StreamAdd    bool `json:"stream_add"`    // Create a stream.
	StreamEdit   bool `json:"stream_edit"`   // Edit a stream.
	StreamRemove bool `json:"stream_remove"` // Remove a stream.

	// Station permissions.
	StationAdd    bool `json:"station_add"`    // Create a station.
	StationEdit   bool `json:"station_edit"`   // Edit a station.
	StationRemove bool `json:"station_remove"` // Remove a station.
}

// Initialize the database tables.
func initIciiTables(d *gorm.DB) {
	d.AutoMigrate(&Station{})
	d.AutoMigrate(&Tag{})
	d.AutoMigrate(&Track{})
	d.AutoMigrate(&UserPermission{})
	d.AutoMigrate(&User{})
}
