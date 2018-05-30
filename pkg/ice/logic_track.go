package ice

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"broadcastle.co/code/icii/pkg/database"
	"github.com/bogem/id3v2"
	"github.com/labstack/echo"
)

// Track information
type Track struct {
	database.Track
}

// Create makes a new track entry.
func (t *Track) Create() error {

	t.FixTags()

	if err := db.Create(&t).Error; err != nil {
		os.Remove(t.Location)
		return err
	}

	return nil

}

// Get the first track with this info.
func (t *Track) Get() error {

	return db.Where(&t).First(&t).Error

}

// Update the track with the data from info.
func (t *Track) Update(i interface{}) error {

	info := i.(Track)

	return db.Model(&t).Updates(info).Error

}

// Delete the track
func (t *Track) Delete() error {

	if err := os.Remove(t.Location); err != nil {
		return err
	}

	return db.Delete(&t).Error

}

// Echo fills in t from data in c.
func (t *Track) Echo(c echo.Context) error {

	i := c.Param("track")

	id, err := strconv.Atoi(i)
	if err != nil {
		return err
	}

	t.ID = uint(id)

	return t.Get()

}

// FixTags fixes any issues with the tags.
func (t *Track) FixTags() {

	tag, err := id3v2.Open(t.Location, id3v2.Options{Parse: true})
	if err != nil {

		log.Println(err)

		if t.Title == "" {
			t.Title = "import from " + filepath.Base(t.Location)
		}

		if t.Artist == "" {
			t.Artist = "Unknown Artist"
		}

		return
	}

	switch {
	case t.Title != "":
	case t.Title == "" && tag.Title() != "":
		t.Title = tag.Title()
	default:
		t.Title = filepath.Base(t.Location)
	}

	if t.Artist == "" {
		t.Artist = tag.Artist()
	}

	if t.Album == "" {
		t.Album = tag.Album()
	}

	if t.Genre == "" {
		t.Genre = tag.Genre()
	}

	if t.Year == "" {
		t.Year = tag.Year()
	}

}
