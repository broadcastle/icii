package web

import (
	"net/http"
	"strconv"

	"broadcastle.co/code/icii/pkg/database"
	"github.com/labstack/echo"
)

func getSong(c echo.Context) error {

	i := c.Param("id")

	id, err := strconv.Atoi(i)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	var song database.Song

	db.First(&song, id)

	if song.ID == 0 {
		return c.JSON(http.StatusNotFound, "not found")
	}

	return c.JSON(http.StatusOK, song)
}

func uploadSong(c echo.Context) error {

	var song database.Song

	if err := c.Bind(&song); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	db.Create(&song)

	return c.JSON(http.StatusOK, song)
}

func updateSong(c echo.Context) error {

	i := c.Param("id")

	id, err := strconv.Atoi(i)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	var original database.Song
	var update database.Song

	db.First(&original, id)

	if original.ID == 0 {
		return c.JSON(http.StatusNotFound, "song does not exist")
	}

	if err := c.Bind(&update); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	db.Model(&original).Updates(update)

	// db.First(&original, id)

	return c.JSON(http.StatusOK, original)

}

func deleteSong(c echo.Context) error {

	i := c.Param("id")

	id, err := strconv.Atoi(i)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	var song database.Song

	db.First(&song, id)

	if song.ID == 0 {
		return c.JSON(http.StatusNotFound, "not found")
	}

	db.Delete(&song)

	return c.JSON(http.StatusOK, "successfully deleted")
}
