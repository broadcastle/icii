package web

import (
	"net/http"

	"github.com/labstack/echo"
)

func userCreate(c echo.Context) error {

	// Bind the sent data to the User struct.
	var user User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// Create the user.
	if err := user.Create(); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	return c.JSON(http.StatusOK, user)
}

func userLogin(c echo.Context) error {

	// Bind the email and password
	var user User

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

func userDelete(c echo.Context) error {

	user, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	if err := user.Delete(); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	return c.JSON(http.StatusOK, "user was deleted")

}

func userUpdate(c echo.Context) error {

	user, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	var new User

	if err := c.Bind(&new); err != nil {
		return c.JSON(msg(http.StatusInternalServerError, err))
	}

	if err := user.Update(new); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	return c.JSON(http.StatusOK, user)

}

func userRetrieve(c echo.Context) error {

	user, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	if err := user.Get(); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	return c.JSON(http.StatusOK, user)

}
