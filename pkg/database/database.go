package database

import (
	"github.com/jinzhu/gorm"
)

// Track is the information needed to store a audio file.
type Track struct {
	gorm.Model

	// Necessary information.
	Album  string `json:"album" gorm:"default:'Unknown Album'"`
	Artist string `json:"artist" gorm:"default:'Unknown Artist'"`
	Title  string `json:"title" gorm:"default:'Unknown Title'"`
	Year   string `json:"year" gorm:"default:'1776'"`

	// Optional information.
	Genre string `json:"genre"`
	Tags  []Tag  `json:"tags"`

	// Calculated information.
	Processed bool    `json:"processed"`
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
	Stations []*Station `gorm:"many2many:user_stations;" json:"stations"`
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

	// Stations
	Stream   Stream
	StreamID uint `json:"stream"`
}

// UserPermission keeps track of what permission are allowed for a user.
type UserPermission struct {
	gorm.Model

	UserID    uint `json:"user_id"`
	StationID uint `json:"station_id"`

	// Station track permissions.
	TrackRead  bool `json:"track_read"`  // See what tracks are available.
	TrackWrite bool `json:"track_write"` // Edit, add, or remove tracks.

	// Station user permissions.
	UserRead  bool `json:"user_read"`  // See the users in the station.
	UserWrite bool `json:"user_write"` // Add or remove users, or edit their permissions.

	// Station stream permissions.
	StreamRead  bool `json:"stream_read"`  // See the station's streams.
	StreamWrite bool `json:"stream_write"` // Edit, add, or remove streams.

	// Station permissions.
	StationRead  bool `json:"station_read"`  // See private station information.
	StationWrite bool `json:"station_write"` // Edit, add, or remove station information.

	// Schedule permissions.
	ScheduleRead  bool `json:"schedule_read"`  // See the schedule for the station.
	ScheduleWrite bool `json:"schedule_write"` // Edit the station schedule.

	// Playlist permissions.
	PlaylistRead  bool `json:"playlist_read"`  // See the playlist.
	PlaylistWrite bool `json:"playlist_write"` // Edit, add, or remove a playlist.

}

// Stream has the information about a stream.
type Stream struct {
	gorm.Model

	Host        string `json:"string"`
	Port        int    `json:"port"`
	Mount       string `json:"mount"`
	User        string `json:"user"`
	Password    string `json:"password"`
	Name        string `json:"stream_name"`
	URL         string `json:"url"`
	Genre       string `json:"genre"`
	Description string `json:"description"`
	BufferSize  int    `json:"buffer_size"`
}

// Initialize the database tables.
func initIciiTables(d *gorm.DB) {
	d.AutoMigrate(&Station{})
	d.AutoMigrate(&Tag{})
	d.AutoMigrate(&Track{})
	d.AutoMigrate(&UserPermission{})
	d.AutoMigrate(&User{})
	d.AutoMigrate(&Stream{})
}
