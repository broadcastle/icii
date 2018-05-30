package ice

import (
	"errors"

	"broadcastle.co/code/icii/pkg/database"
	slugify "github.com/mozillazg/go-slugify"
)

// Station info
type Station struct {
	database.Station
}

// Create s
func (s *Station) Create() error {

	if s.Name == "" {
		return errors.New("need a station name")
	}

	if s.Slug == "" {
		s.Slug = slugify.Slugify(s.Name)
	}

	var found Station
	if err := db.Where("slug = ?", s.Slug).First(&found).Error; err == nil {
		return errors.New("station exists")
	}

	return db.Create(&s).Error
}

// Update s with info
func (s *Station) Update(i interface{}) error {

	info := i.(Station)

	if s.Slug == "" {
		s.Slug = slugify.Slugify(s.Name)
	}

	var found Station
	if err := db.Where("slug = ?", s.Slug).First(&found).Error; err == nil {
		return errors.New("station is using this slug")
	}

	return db.Model(&s).Updates(info).Error
}

// Delete the station.
func (s *Station) Delete() error {
	return db.Delete(&s).Error
}

// Get the station information.
func (s *Station) Get() error {
	return db.Where(&s).First(&s).Error
}
