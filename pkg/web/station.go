package web

import (
	"net/http"

	"broadcastle.co/code/icii/pkg/database"
	"github.com/labstack/echo"
	slugify "github.com/mozillazg/go-slugify"
)

func stationCreate(c echo.Context) error {

	// Get the user ID if the token is valid.
	userID, err := getJwtID(c)
	if err != nil {
		return c.JSON(msg(http.StatusInternalServerError, err))
	}

	//// Bind the station information.
	var org database.Station

	if err := c.Bind(&org); err != nil {
		return c.JSON(msg(http.StatusInternalServerError, err))
	}

	if org.Name == "" {
		return c.JSON(msg(http.StatusMethodNotAllowed, "need a station name"))
	}

	if org.Slug == "" {
		org.Slug = slugify.Slugify(org.Name)
	}

	// Save the database.
	if err := db.Create(&org).Error; err != nil {
		return c.JSON(msg(http.StatusInternalServerError, err))
	}

	//// Create and save the user permissions for the station.
	permissions := database.UserPermission{
		StationID:     org.ID,
		UserID:        userID,
		TrackRead:     true,
		TrackWrite:    true,
		UserRead:      true,
		UserWrite:     true,
		StreamRead:    true,
		StreamWrite:   true,
		StationRead:   true,
		StationWrite:  true,
		ScheduleRead:  true,
		ScheduleWrite: true,
		PlaylistRead:  true,
		PlaylistWrite: true,
	}

	if err := db.Create(&permissions).Error; err != nil {
		return c.JSON(msg(http.StatusInternalServerError, err))
	}

	return c.JSON(msg(http.StatusOK, org))
}
