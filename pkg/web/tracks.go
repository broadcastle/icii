package web

import (
	"net/http"

	"broadcastle.co/code/icii/pkg/ice"
	"github.com/labstack/echo"
)

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

	var track ice.Track

	track.UserID = userID
	track.StationID = stationID
	track.Title = c.FormValue("title")
	track.Album = c.FormValue("album")
	track.Artist = c.FormValue("artist")
	track.Year = c.FormValue("year")
	track.Genre = c.FormValue("genre")

	// Get the audio file
	file, err := c.FormFile("audio")
	if err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	if err := track.Upload(file); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	return c.JSON(http.StatusOK, "track is being processed")
}

// Retrieve a track when given an ID.
func trackGet(c echo.Context) error {

	// Check if the token is valid.
	if _, err := getJwtID(c); err != nil {
		return c.JSON(msg(http.StatusForbidden, err))
	}

	// Get the ID as an integer.
	id, err := getIDfromParam("track", c)
	if err != nil {
		return c.JSON(msg(http.StatusInternalServerError, err))
	}

	// Find the track with that ID and return the data.
	var track ice.Track
	track.ID = id

	if err := track.Get(); err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, track)
}

// Update the track at the given ID.
func trackUpdate(c echo.Context) error {

	// Check if the token is valid.
	if _, err := getJwtID(c); err != nil {
		return c.JSON(msg(http.StatusForbidden, err))
	}

	// //// Get the ID as an integer
	// id, err := getIDfromParam("track", c)
	// if err != nil {
	// 	return c.JSON(msg(http.StatusInternalServerError, err))
	// }

	// Bind the updated information to a Track struct.
	var update ice.Track
	if err := c.Bind(&update); err != nil {
		return c.JSON(msg(http.StatusInternalServerError, err))
	}

	// var original ice.Track
	// original.ID = id

	// if err := original.Get(); err != nil {
	// 	return c.JSON(http.StatusNotFound, err)
	// }

	// if err := original.Update(update); err != nil {
	// 	return c.JSON(http.StatusMethodNotAllowed, err)
	// }

	track := ice.InitTrack()

	if err := ice.Echo(track, c); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	if err := ice.Update(track, update); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	return c.JSON(http.StatusOK, track)

}

// Delete the track at the given ID.
func trackDelete(c echo.Context) error {

	// Check if the token is valid.
	if _, err := getJwtID(c); err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	// Get the ID as an iteger and check that it's not 0.

	id, err := getIDfromParam("track", c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	var track ice.Track
	track.ID = id

	if err := track.Get(); err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	if err := track.Delete(); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	return c.JSON(http.StatusOK, "successfully deleted")
}
