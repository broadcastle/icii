package database

import (
	"github.com/jinzhu/gorm"
)

// Track is the information needed to store a audio file.
type Track struct {
	gorm.Model
	Album          string  `json:"album"`
	Artist         string  `json:"artist"`
	Location       string  `json:"location"`
	Title          string  `json:"title"`
	OrganizationID uint    `json:"organization_id"`
	UserID         uint    `json:"uploader"`
	Year           string  `json:"year"`
	Length         float64 `json:"length"`
	Genre          string  `json:"genre"`
	Tags           []Tag   `json:"tags"`
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
	Name          string          `json:"name"`
	Email         string          `json:"email"`
	Password      string          `json:"password"`
	Organizations []*Organization `gorm:"many2many:user_organizations;" json:"organizations"`
}

// A Organization holds the audio tracks and users.
type Organization struct {
	gorm.Model
	Name  string  `json:"name"`
	Slug  string  `json:"slug"`
	Users []*User `gorm:"many2many:user_organizations;" json:"users"`
	Track []Track `json:"audio"`
}

// UserPermission keeps track of what permission are allowed for a user.
type UserPermission struct {
	gorm.Model
	UserID         uint
	OrganizationID uint

	TrackAdd    bool
	TrackEdit   bool
	TrackRemove bool

	UserAdd    bool
	UserEdit   bool
	UserRemove bool

	StreamAdd    bool
	StreamEdit   bool
	StreamRemove bool

	OrganizationAdd    bool
	OrganizationEdit   bool
	OrganizationRemove bool
}

// Initialize the database tables.
func initIciiTables(d *gorm.DB) {
	d.AutoMigrate(&Organization{})
	d.AutoMigrate(&Tag{})
	d.AutoMigrate(&Track{})
	d.AutoMigrate(&User{})
	d.AutoMigrate(&UserPermission{})
}
