package web

import (
	"net/http"

	"broadcastle.co/code/icii/pkg/database"
	"github.com/labstack/echo"
	slugify "github.com/mozillazg/go-slugify"
)

func stationCreate(c echo.Context) error {

	var org database.Station

	c.Bind(&org)

	if org.Name == "" {
		return c.JSON(http.StatusMethodNotAllowed, "need a station name")
	}

	if org.Slug == "" {
		org.Slug = slugify.Slugify(org.Name)
	}

	db.Create(&org)

	userID := getJwtID(c)

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

	db.Create(&permissions)

	return c.JSON(http.StatusOK, org)
}
