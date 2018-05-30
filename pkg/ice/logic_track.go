package ice

import (
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"

	"broadcastle.co/code/icii/pkg/database"
	"github.com/bogem/id3v2"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	filetype "gopkg.in/h2non/filetype.v1"
)

// Track information
type Track struct {
	database.Track
}

// Upload will upload a file with with track information.
func (t *Track) Upload(file *multipart.FileHeader) error {

	src, err := file.Open()
	if err != nil {
		return err
	}

	u := uuid.NewV4()

	ext := path.Ext(file.Filename)
	tmp := path.Join(os.TempDir(), u.String()+ext)

	dst, err := os.Create(tmp)
	if err != nil {
		return err
	}

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	t.Location = tmp

	// Process the track. Let the user know the track is being processed.

	go t.ProcessFile()

	return nil
}

// ProcessFile edits the track.
func (t *Track) ProcessFile() {

	defer os.Remove(t.Location)

	in, err := os.Open(t.Location)
	if err != nil {
		log.Println(err)
		return
	}

	head := make([]byte, 261)
	if _, err := in.Read(head); err != nil {
		log.Println(err)
		return
	}

	if !filetype.IsMIME(head, "audio/mpeg") {
		os.Remove(t.Location)
		return
	}

	end := viper.GetString("files.location")
	_, filename := path.Split(t.Location)
	end = path.Join(end, filename)

	out, err := os.Create(end)
	if err != nil {
		log.Println(err)
		return
	}

	if _, err := io.Copy(out, in); err != nil {
		log.Println(err)
		return
	}

	in.Close()
	out.Close()

	t.Location = end

	// Process the Tags
	t.FixTags()

	// Create the database entry

	if err := db.Create(&t).Error; err != nil {
		log.Println(err)
		os.Remove(end)
	}

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

// Get the first track with this info.
func (t *Track) Get() error {

	return db.Where(&t).First(&t).Error

}

// Update the track with the data from info.
func (t *Track) Update(info Track) error {

	return db.Model(&t).Updates(info).Error

}

// Delete the track
func (t *Track) Delete() error {

	if err := os.Remove(t.Location); err != nil {
		return err
	}

	return db.Delete(&t).Error
}
