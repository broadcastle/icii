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
func processTrack(location string, info database.Track, filename string) {

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

// Create a track entry in the database.
func trackCreate(c echo.Context) error {

	// Check if the token is valid.
	userID, err := getJwtID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	//// Bind the sent data to a entry.
	var track database.Track

	track.Title = c.FormValue("title")
	track.Album = c.FormValue("album")
	track.Artist = c.FormValue("artist")
	track.Year = c.FormValue("year")
	track.Genre = c.FormValue("genre")
	track.UserID = userID

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

	//// Create the destination
	dst, err := os.Create(tmp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if _, err := io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	//// Process the track. Let the user know the track is being processed.
	go processTrack(tmp, track, file.Filename)

	return c.JSON(http.StatusOK, "track is being processed")
}

// Retrieve a track when given an ID.
func trackGet(c echo.Context) error {

	// Check if the token is valid.
	if _, err := getJwtID(c); err != nil {
		return c.JSON(http.StatusInternalServerError, msg(err))
	}

	//// Get the ID as an integer.
	i := c.Param("track")

	id, err := strconv.Atoi(i)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, msg(err))
	}

	//// Find the track with that ID and return the data.
	var track database.Track

	if err := db.First(&track, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, msg(err))
	}

	if track.ID == 0 {
		return c.JSON(http.StatusNotFound, "not found")
	}

	return c.JSON(http.StatusOK, track)
}

// Update the track at the given ID.
func trackUpdate(c echo.Context) error {

	// Check if the token is valid.
	if _, err := getJwtID(c); err != nil {
		return c.JSON(http.StatusInternalServerError, msg(err))
	}

	//// Get the ID as an integer
	i := c.Param("track")

	id, err := strconv.Atoi(i)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, msg(err))
	}

	//// Get the original track info.
	var original database.Track

	if err := db.First(&original, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, msg(err))
	}

	//// Bind the updated information to a Track struct.
	var update database.Track
	if err := c.Bind(&update); err != nil {
		return c.JSON(http.StatusInternalServerError, msg(err))
	}

	//// Update the original track and return the result.
	if err := db.Model(&original).Updates(update).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, msg(err))
	}

	return c.JSON(http.StatusOK, original)

}

// Delete the track at the given ID.
func trackDelete(c echo.Context) error {

	// Check if the token is valid.
	if _, err := getJwtID(c); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// Get the ID as an iteger and check that it's not 0.
	i := c.Param("track")

	id, err := strconv.Atoi(i)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// if id == 0 {
	// 	return c.JSON(http.StatusMethodNotAllowed, "an id is needed")
	// }

	// Query the first track that has the ID and delete it.
	var track database.Track

	if err := db.First(&track, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err)

	}

	if err := db.Delete(&track).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "successfully deleted")
}
