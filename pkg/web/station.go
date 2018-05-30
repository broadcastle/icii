package web

import (
	"net/http"

	"broadcastle.co/code/icii/pkg/ice"
	"github.com/labstack/echo"
)

func stationCreate(c echo.Context) error {

	user, err := ice.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	// Bind the station information.
	var info ice.Station

	if err := c.Bind(&info); err != nil {
		return c.JSON(msg(http.StatusInternalServerError, err))
	}

	// Create the station.
	if err := user.CreateStation(info); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	return c.JSON(http.StatusOK, info)
}
