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

	if err := db.Where("email = ?", user.Email).First(&found).Error; err == nil {
		return c.JSON(http.StatusInternalServerError, msg(err))
	}

	//// Encrypt the password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	user.Password = string(hash)

	//// Create the user and return the result.

	if err := db.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, msg(err))
	}

	return c.JSON(http.StatusOK, user)
}

func userLogin(c echo.Context) error {

	//// Bind the email and password
	var sent database.User

	wrong := "email and/or password was incorrect"

	if err := c.Bind(&sent); err != nil {
		return c.JSON(http.StatusInternalServerError, msg(err))
	}

	if sent.Email == "" && sent.Password == "" {
		return c.JSON(http.StatusMethodNotAllowed, msg(wrong))
	}

	//// Find the user
	var user database.User

	if err := db.Where("email = ?", sent.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, msg(err))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(sent.Password)); err != nil {
		return c.JSON(http.StatusMethodNotAllowed, msg(wrong))
	}

	//// Create the JWT token.

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(viper.GetString("icii.jwt")))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, msg(err))
	}

	return c.JSON(http.StatusOK, msg(t))
}

func userDelete(c echo.Context) error {

	// Get user ID if the token is valid.
	userID, err := getJwtID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	//// Delete the user
	var user database.User

	user.ID = userID

	if err := db.Delete(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, msg(err))
	}

	return c.JSON(http.StatusOK, msg("user was deleted"))

}

func userUpdate(c echo.Context) error {

	// Get the current user information.
	userID, err := getJwtID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	var user database.User

	if err := db.First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, msg(err))
	}

	// Get updated user information.
	var new database.User

	if err := c.Bind(&new); err != nil {
		return c.JSON(http.StatusInternalServerError, msg(err))
	}

	// If there is a password, encrypt it.
	if new.Password != "" {

		hash, err := bcrypt.GenerateFromPassword([]byte(new.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, msg(err))
		}

		new.Password = string(hash)

	}

	if err := db.Model(&user).Updates(new).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, msg(err))
	}

	return c.JSON(http.StatusOK, user)

}

func userRetrieve(c echo.Context) error {

	// Get the user ID if the token is valid.
	userID, err := getJwtID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// Get the user information.
	var user database.User

	if err := db.Find(&user, userID).Error; err != nil {
		return c.JSON(http.StatusMethodNotAllowed, msg(err))
	}

	return c.JSON(http.StatusOK, user)

}
