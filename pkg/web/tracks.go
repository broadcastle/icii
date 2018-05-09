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
	filetype "gopkg.in/h2non/filetype.v1"
)

func processTags(name, location string, info database.Track) database.Track {

	info.Location = location

	tag, err := id3v2.Open(location, id3v2.Options{Parse: true})
	if err != nil {

		log.Printf("error processing file: %s\n%v\n", location, err)

		if info.Title == "" {
			info.Title = "imported from " + name
		}

		if info.Artist == "" {
			info.Artist = "Unknown Artist"
		}
		return info

	}

	switch {
	case info.Title != "":
	case info.Title == "" && tag.Title() != "":
		info.Title = tag.Title()
	default:
		info.Title = name

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

	return info
}

// Process the audio file that was uploaded.
func processTrack(location string, info database.Track, originalName string) {

	defer os.Remove(location)

	// Open the file
	in, err := os.Open(location)
	if err != nil {
		log.Println(err)
		return
	}

	defer in.Close()

	// Get the header
	head := make([]byte, 261)
	if _, err := in.Read(head); err != nil {
		log.Println(err)
		return
	}

	if _, err := in.Seek(0, 0); err != nil {
		log.Println(err)
		return
	}

	// Remove file if it is not a mp3 file.
	if !filetype.IsMIME(head, "audio/mpeg") {
		os.Remove(location)
		return
	}

	// Get the location to store the files.
	end := viper.GetString("files.location")
	_, filename := path.Split(location)
	end = path.Join(end, filename)

	// Move the file to the new location
	out, err := os.Create(end)
	if err != nil {
		log.Println(err)
		return
	}

	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		log.Println(err)
		return
	}

	// Process the Tags
	info = processTags(originalName, end, info)

	// Create the database entry

	if err := db.Create(&info).Error; err != nil {
		log.Println(err)
		os.Remove(end)
	}
}

// Create a track entry in the database.
func trackCreate(c echo.Context) error {

	// Check if the token is valid.
	userID, err := getJwtID(c)
	if err != nil {
		return c.JSON(msg(http.StatusMethodNotAllowed, err))
	}

	//// Bind the sent data to a entry.

	stationID, err := getStationID(c.FormValue("station"))
	if err != nil {
		return c.JSON(msg(http.StatusMethodNotAllowed, err))
	}

	track := database.Track{
		UserID:    userID,
		StationID: stationID,
		Title:     c.FormValue("title"),
		Album:     c.FormValue("album"),
		Artist:    c.FormValue("artist"),
		Year:      c.FormValue("year"),
		Genre:     c.FormValue("genre"),
	}

	//// Copy the audio file to a temporary folder

	// Get the audio file
	file, err := c.FormFile("audio")
	if err != nil {
		return c.JSON(msg(http.StatusMethodNotAllowed, err))
	}

	// Source
	src, err := file.Open()
	if err != nil {
		return c.JSON(msg(http.StatusInternalServerError, err))
	}

	//// Prepare a new name for the file.

	// Generate UUID
	u := uuid.NewV4()

	//// Create a temporary destination
	ext := path.Ext(file.Filename)
	tmp := path.Join(os.TempDir(), u.String()+ext)

	dst, err := os.Create(tmp)
	if err != nil {
		return c.JSON(msg(http.StatusInternalServerError, err))
	}

	if _, err := io.Copy(dst, src); err != nil {
		return c.JSON(msg(http.StatusInternalServerError, err))
	}

	//// Process the track. Let the user know the track is being processed.
	go processTrack(tmp, track, file.Filename)

	return c.JSON(http.StatusOK, "track is being processed")
}

// Retrieve a track when given an ID.
func trackGet(c echo.Context) error {

	// Check if the token is valid.
	if _, err := getJwtID(c); err != nil {
		return c.JSON(msg(http.StatusForbidden, err))
	}

	//// Get the ID as an integer.
	i := c.Param("track")

	id, err := strconv.Atoi(i)
	if err != nil {
		return c.JSON(msg(http.StatusInternalServerError, err))
	}

	//// Find the track with that ID and return the data.
	var track database.Track

	if err := db.First(&track, id).Error; err != nil {
		return c.JSON(msg(http.StatusInternalServerError, err))
	}

	if track.ID == 0 {
		return c.JSON(msg(http.StatusNotFound, "not found"))
	}

	return c.JSON(msg(http.StatusOK, track))
}

// Update the track at the given ID.
func trackUpdate(c echo.Context) error {

	// Check if the token is valid.
	if _, err := getJwtID(c); err != nil {
		return c.JSON(msg(http.StatusForbidden, err))
	}

	//// Get the ID as an integer
	i := c.Param("track")

	id, err := strconv.Atoi(i)
	if err != nil {
		return c.JSON(msg(http.StatusInternalServerError, err))
	}

	//// Get the original track info.
	var original database.Track

	if err := db.First(&original, id).Error; err != nil {
		return c.JSON(msg(http.StatusInternalServerError, err))
	}

	//// Bind the updated information to a Track struct.
	var update database.Track
	if err := c.Bind(&update); err != nil {
		return c.JSON(msg(http.StatusInternalServerError, err))
	}

	//// Update the original track and return the result.
	if err := db.Model(&original).Updates(update).Error; err != nil {
		return c.JSON(msg(http.StatusInternalServerError, err))
	}

	return c.JSON(msg(http.StatusOK, original))

}

// Delete the track at the given ID.
func trackDelete(c echo.Context) error {

	// Check if the token is valid.
	if _, err := getJwtID(c); err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	// Get the ID as an iteger and check that it's not 0.
	i := c.Param("track")

	id, err := strconv.Atoi(i)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

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
