package web

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
func (s *Station) Update(info Station) error {

	if s.Slug == "" {
		s.Slug = slugify.Slugify(s.Name)
	}

	var found Station
	if err := db.Where("slug = ?", s.Slug).First(&found).Error; err == nil {
		return errors.New("station is using this slug")
	}

	return db.Model(&s).Updates(info).Error
}

// Delete s
func (s *Station) Delete() error {
	return errors.New("need to be written")
}
