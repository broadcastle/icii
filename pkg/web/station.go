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

func stationDelete(c echo.Context) error {

	station := ice.InitStation()

	if err := ice.Echo(c, station); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	if err := ice.Remove(station); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	return c.JSON(http.StatusOK, "sucess")

}

func stationUpdate(c echo.Context) error {

	station := ice.InitStation()

	if err := ice.Echo(c, station); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	var new ice.Station

	if err := c.Bind(&new); err != nil {
		return c.JSON(msg(http.StatusInternalServerError, err))
	}

	if err := ice.Update(station, new); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	return c.JSON(http.StatusOK, station)

}

func stationRetrieve(c echo.Context) error {

	station := ice.InitStation()

	if err := ice.Echo(c, station); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	return c.JSON(http.StatusOK, station)

}
