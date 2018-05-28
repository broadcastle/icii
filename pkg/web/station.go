package web

import (
	"net/http"

	"github.com/labstack/echo"
)

func stationCreate(c echo.Context) error {

	user, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	// Bind the station information.
	var info Station

	if err := c.Bind(&info); err != nil {
		return c.JSON(msg(http.StatusInternalServerError, err))
	}

	// Create the station.
	if err := user.CreateStation(info); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	return c.JSON(http.StatusOK, info)
}
