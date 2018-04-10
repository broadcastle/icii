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
