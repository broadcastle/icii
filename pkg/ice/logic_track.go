package ice

import (
	"os"
	"path/filepath"
	"strconv"

	"broadcastle.co/code/icii/pkg/database"
	"github.com/bogem/id3v2"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
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
		log.Warn(err)
		return err
	}

	log.Infof("icii created track #%x for user #%x", t.ID, t.UserID)

	return nil

}

// Get the first track with this info.
func (t *Track) Get() error {

	if err := db.Where(&t).First(&t).Error; err != nil {
		log.Warn(err)
		return err
	}

	log.Infof("icii retrieved information for track #%x", t.ID)

	return nil

}

// Update the track with the data from info.
func (t *Track) Update(i interface{}) error {

	info := i.(Track)

	if err := db.Model(&t).Updates(info).Error; err != nil {
		log.Warn(err)
		return err
	}

	log.Infof("icii updated track #%x with new information", t.ID)

	return nil

}

// Delete the track
func (t *Track) Delete() error {

	if err := db.Delete(&t).Error; err != nil {
		log.Warn(err)
		return err
	}

	if err := os.Remove(t.Location); err != nil {
		log.Warn("err")
		return err
	}

	log.Infof("icii removed track #%x", t.ID)

	return nil

}

// Echo fills in t from data in c.
func (t *Track) Echo(c echo.Context) error {

	i := c.Param("track")

	id, err := strconv.Atoi(i)
	if err != nil {
		log.Warn(err)
		return err
	}

	s := c.Param("station")
	sid, err := strconv.Atoi(s)
	if err != nil {
		log.Warn(err)
		return err
	}

	t.ID = uint(id)
	t.StationID = uint(sid)

	return t.Get()

}

// FixTags fixes any issues with the tags.
func (t *Track) FixTags() {

	tag, err := id3v2.Open(t.Location, id3v2.Options{Parse: true})
	if err != nil {

		log.Warnf("icii was unable to read id3v2 tags for track at %s: %x", t.Location, err)

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

// Play will take a valid t and play it.
func (t Track) Play() error {
	return nil
}
