package web

import (
	"net/http"

	"broadcastle.co/code/icii/pkg/ice"
	"github.com/labstack/echo"
)

func userCreate(c echo.Context) error {

	// Bind the sent data.
	user := ice.InitUser()

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// Create the user.
	if err := ice.New(user); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	return c.JSON(http.StatusOK, user)
}

func userDelete(c echo.Context) error {

	user := ice.InitUser()

	if err := ice.Echo(c, user); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	if err := ice.Remove(user); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	return c.JSON(http.StatusOK, "user was deleted")

}

func userUpdate(c echo.Context) error {

	user := ice.InitUser()

	if err := ice.Echo(c, user); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	var new ice.User

	if err := c.Bind(&new); err != nil {
		return c.JSON(msg(http.StatusInternalServerError, err))
	}

	if err := ice.Update(user, new); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	return c.JSON(http.StatusOK, user)

}

func userRetrieve(c echo.Context) error {

	user := ice.InitUser()

	if err := ice.Echo(c, user); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	return c.JSON(http.StatusOK, user)

}

func userLogin(c echo.Context) error {

	// Bind the email and password
	var user ice.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// Login
	t, err := user.Login()
	if err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	return c.JSON(msg(http.StatusOK, t))
}
