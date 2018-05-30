package web

import (
	"net/http"

	"broadcastle.co/code/icii/pkg/ice"
	"github.com/labstack/echo"
)

func stationCreate(c echo.Context) error {

	station := ice.InitStation()

	if err := c.Bind(&station); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if err := ice.New(station); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	// CLEAN UP
	user, err := ice.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if err := user.(*ice.User).CreateStation(*station.(*ice.Station)); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	return c.JSON(http.StatusOK, station)

}

func stationUpdate(c echo.Context) error {
	return nil
}

func stationRetrieve(c echo.Context) error {
	return nil
}

func stationDelete(c echo.Context) error {
	return nil
}
