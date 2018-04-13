package web

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// This is a fast way to get the user id.
func getJwtID(c echo.Context) uint {

	i := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["id"].(float64)

	return uint(i)

}
