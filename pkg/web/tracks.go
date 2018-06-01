package web

import (
	"net/http"

	"broadcastle.co/code/icii/pkg/ice"
	"github.com/labstack/echo"
)

// Create a track entry in the database.
func trackCreate(c echo.Context) error {

	user := ice.InitUser()

	if err := ice.Echo(user, c); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	station := ice.InitStation()

	if err := ice.Echo(station, c); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	file, err := c.FormFile("audio")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	location, err := ice.FormImportFile(file)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	track := ice.InitTrack()

	track.(*ice.Track).UserID = user.(*ice.User).ID
	track.(*ice.Track).Title = c.FormValue("title")
	track.(*ice.Track).Album = c.FormValue("album")
	track.(*ice.Track).Artist = c.FormValue("artist")
	track.(*ice.Track).Year = c.FormValue("year")
	track.(*ice.Track).Genre = c.FormValue("genre")
	track.(*ice.Track).Location = location
	track.(*ice.Track).StationID = station.(*ice.Station).ID

	if err := ice.New(track); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	return c.JSON(http.StatusOK, "success")

}

// Retrieve a track when given an ID.
func trackGet(c echo.Context) error {

	// Check if the token is valid.
	if _, err := getJwtID(c); err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	track := ice.InitTrack()

	if err := ice.Echo(track, c); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, track)
}

// Update the track at the given ID.
func trackUpdate(c echo.Context) error {

	// Check if the token is valid.
	if _, err := getJwtID(c); err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	// Bind the updated information to a Track struct.
	var update ice.Track
	if err := c.Bind(&update); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

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

	track := ice.InitTrack()

	if err := ice.Echo(track, c); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if err := ice.Remove(track); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	return c.JSON(http.StatusOK, "success")
}

func trackPlay(c echo.Context) error {

	track := ice.InitTrack()

	if err := ice.Echo(track, c); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	if err := track.(*ice.Track).Play(); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "playing")
}
