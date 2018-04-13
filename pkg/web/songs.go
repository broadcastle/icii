package web

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"broadcastle.co/code/icii/pkg/database"
	"github.com/bogem/id3v2"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
)

// Process the audio file that was uploaded.
func processSong(location string, info database.Song, filename string) {

	//// Get the tags from the temporary file.

	tag, err := id3v2.Open(location, id3v2.Options{Parse: true})
	if err != nil {

		log.Println(err)

		os.Remove(location)

		return
	}

	//// The title can not be empty. The others can though.
	switch {
	case info.Title != "":
	case info.Title == "" && tag.Title() != "":
		info.Title = tag.Title()
	default:
		info.Title = filename

	}

	if info.Artist == "" {
		info.Artist = tag.Artist()
	}

	if info.Album == "" {
		info.Album = tag.Album()
	}

	if info.Genre == "" {
		info.Genre = tag.Genre()
	}

	if info.Year == "" {
		info.Year = tag.Year()
	}

	info.Location = location

	//// Create the database entry

	db.Create(&info)
}

// Create a song entry in the database.
func songCreate(c echo.Context) error {

	//// Bind the sent data to a entry.
	var song database.Song

	song.Title = c.FormValue("title")
	song.Album = c.FormValue("album")
	song.Artist = c.FormValue("artist")
	song.Year = c.FormValue("year")
	song.Genre = c.FormValue("genre")
	song.UserID = getJwtID(c)

	//// Copy the audio file to a temporary folder

	// Get the audio file
	file, err := c.FormFile("audio")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// Source
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	//// Prepare a new name for the file.

	// Generate UUID
	u := uuid.NewV4()

	// Name the destination
	t := viper.GetString("files.temporary")
	ext := path.Ext(file.Filename)

	tmp := path.Join(t, u.String()+ext)

	log.Println(tmp)

	//// Create the destination
	dst, err := os.Create(tmp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if _, err := io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	//// Process the song. Let the user know the song is being processed.
	go processSong(tmp, song, file.Filename)

	return c.JSON(http.StatusOK, "song is being processed")
}

// Retrieve a song when given an ID.
func songGet(c echo.Context) error {

	//// Get the ID as an integer.
	i := c.Param("id")

	id, err := strconv.Atoi(i)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	//// Find the song with that ID and return the data.
	var song database.Song

	db.First(&song, id)

	if song.ID == 0 {
		return c.JSON(http.StatusNotFound, "not found")
	}

	return c.JSON(http.StatusOK, song)
}

// Update the song at the given ID.
func songUpdate(c echo.Context) error {

	//// Get the ID as an integer
	i := c.Param("id")

	id, err := strconv.Atoi(i)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	//// Get the original song info.
	var original database.Song

	db.First(&original, id)

	if original.ID == 0 {
		return c.JSON(http.StatusNotFound, "song does not exist")
	}

	//// Bind the updated information to a Song struct.
	var update database.Song
	if err := c.Bind(&update); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	//// Update the original song and return the result.
	db.Model(&original).Updates(update)

	return c.JSON(http.StatusOK, original)

}

// Delete the song at the given ID.
func songDelete(c echo.Context) error {

	// Get the ID as an iteger and check that it's not 0.
	i := c.Param("id")

	id, err := strconv.Atoi(i)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if id == 0 {
		return c.JSON(http.StatusMethodNotAllowed, "an id is needed")
	}

	// Query the first song that has the ID and delete it.
	var song database.Song

	db.First(&song, id)

	if song.ID == 0 {
		return c.JSON(http.StatusNotFound, "not found")
	}

	db.Delete(&song)

	return c.JSON(http.StatusOK, "successfully deleted")
}
