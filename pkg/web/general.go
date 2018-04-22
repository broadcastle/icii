package web

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// JSONResponse is used to create a json response.
type JSONResponse struct {
	Msg string `json:"msg"`
}

// This is a fast way to get the user id.
func getJwtID(c echo.Context) uint {

	i := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["id"].(float64)

	return uint(i)

}

func msg(t interface{}) interface{} {

	switch t.(type) {
	case string:
		return &JSONResponse{Msg: t.(string)}
	default:
		return t
	}
}
