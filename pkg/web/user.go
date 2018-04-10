package web

import (
	"net/http"

	"broadcastle.co/code/icii/pkg/database"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func userCreate(c echo.Context) error {

	//// Bind the sent data to the User struct.
	var user database.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	//// Find if a user is already registed with that email.
	var found database.User

	db.Where("email = ?", user.Email).First(&found)

	if found.ID != 0 {
		return c.JSON(http.StatusMethodNotAllowed, "email registered")
	}

	//// Encrypt the password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	user.Password = string(hash)

	//// Create the user and return the result.

	db.Create(&user)

	return c.JSON(http.StatusOK, user)
}
