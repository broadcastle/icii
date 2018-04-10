package web

import (
	"net/http"
	"strconv"

	"broadcastle.co/code/icii/pkg/database"
	"github.com/labstack/echo"
)

// Create a song entry in the database.
func songCreate(c echo.Context) error {

	//// Bind the sent data to a entry.
	var song database.Song

	if err := c.Bind(&song); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	//// Create the database and return the result.
	db.Create(&song)

	return c.JSON(http.StatusOK, song)
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
