package web

import (
	"net/http"
	"time"

	"broadcastle.co/code/icii/pkg/database"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func userCreate(c echo.Context) error {

	//// Bind the sent data to the User struct.
	var user database.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	//// Find if a user is already registered with that email.
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

func userLogin(c echo.Context) error {

	//// Bind the email and password
	var sent database.User

	wrong := "email and/or password was incorrect"

	if err := c.Bind(&sent); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if sent.Email == "" && sent.Password == "" {
		return c.JSON(http.StatusMethodNotAllowed, wrong)
	}

	//// Find the user
	var user database.User

	db.Where("email = ?", sent.Email).First(&user)

	if user.ID == 0 {
		return c.JSON(http.StatusMethodNotAllowed, wrong)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(sent.Password)); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, wrong)
	}

	//// Create the JWT token.

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user.Name
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(viper.GetString("icii.jwt")))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, t)
}
